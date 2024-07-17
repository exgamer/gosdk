package consumer

import (
	"context"
	"errors"
	"github.com/exgamer/gosdk/pkg/config"
	"github.com/exgamer/gosdk/pkg/logger"
	"github.com/exgamer/gosdk/pkg/sentry"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)
import "github.com/IBM/sarama"

func RunConsumer(baseConfig *config.BaseConfig, brokerList []string, handlers map[string]IKafkaHandler) {
	keepRunning := true
	log.Println("Starting a new Sarama consumer")

	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	version, err := sarama.ParseKafkaVersion(sarama.DefaultVersion.String())

	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Version = version

	//switch assignor {
	//case "sticky":
	//	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	//case "roundrobin":
	//	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	//case "range":
	//	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	//default:
	//	log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	//}

	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer := NewKafkaConsumer(make(chan bool), handlers, baseConfig)
	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokerList, baseConfig.Name, config)

	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			topics := make([]string, len(handlers))
			i := 0

			for k := range handlers {
				topics[i] = k
				i++
			}

			if err := client.Consume(ctx, topics, consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}

			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			if consumptionIsPaused {
				client.ResumeAll()
				log.Println("Resuming consumption")
			} else {
				client.PauseAll()
				log.Println("Pausing consumption")
			}

			consumptionIsPaused = !consumptionIsPaused
		}
	}

	cancel()
	wg.Wait()

	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func NewKafkaConsumer(ready chan bool, handlers map[string]IKafkaHandler, baseConfig *config.BaseConfig) *KafkaConsumer {
	return &KafkaConsumer{
		Ready:      ready,
		handlers:   handlers,
		baseConfig: baseConfig,
	}
}

// KafkaConsumer represents a Sarama consumer group consumer
type KafkaConsumer struct {
	Ready      chan bool
	handlers   map[string]IKafkaHandler
	baseConfig *config.BaseConfig
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *KafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.Ready)

	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")

				return nil
			}

			topic := message.Topic
			handler := consumer.handlers[topic].(IKafkaHandler)
			msg := "groupId:" + consumer.baseConfig.Name + ";key:" + string(message.Key) + "; value:" + string(message.Value)
			logger.FormattedInfo(consumer.baseConfig.Name, "consumer", message.Topic, 0, "", msg)
			err := handler.Handle(message)

			if err != nil {
				msg = msg + "; error_text:" + err.Error()
				logger.FormattedError(consumer.baseConfig.Name, "consumer", message.Topic, 0, "", msg)
				sentry.SendError("Kafka Consumer Error: "+err.Error(),
					map[string]string{
						"service_name": consumer.baseConfig.Name,
						"env":          consumer.baseConfig.AppEnv,
						"kafka_group":  consumer.baseConfig.Name,
					},
					map[string]interface{}{
						"key":             string(message.Key),
						"value":           string(message.Value),
						"topic_partition": message.Topic,
						"timestamp":       message.Timestamp,
					},
				)
				session.MarkMessage(message, "")
			}

			//log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
