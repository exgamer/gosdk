package producer

import (
	"github.com/IBM/sarama"
)

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)

	return &KafkaProducer{
		producer: producer,
	}, err
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func (p *KafkaProducer) SendMessage(topic string, message string) error {
	_, _, err := p.producer.SendMessage(p.prepareKafkaMessage(topic, message))

	if err != nil {
		return err
	}

	return nil
}

func (p *KafkaProducer) prepareKafkaMessage(topic, message string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}

	return msg
}
