package consumer

import (
	"github.com/IBM/sarama"
)

type IKafkaHandler interface {
	Handle(message *sarama.ConsumerMessage) error
}
