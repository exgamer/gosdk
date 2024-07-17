Пример использования функции для отправки смс через кафку

```go
package app


func main() {
	kafkaProducer, err := NewKafkaProducer(brokers)

	messageSender := NewMessageSender(kafkaProducer)

	err := messageSender.SendSms("+7 (778)-300-00-44", "test message")

	if err != nil {

		println(err)
	}
}

```


Пример использования функции для отправки email через кафку

```go
package app

func main() {
	kafkaProducer, err := NewKafkaProducer(brokers)

	messageSender := NewMessageSender(kafkaProducer)

	err := messageSender.SendEmailWithPlainText("example@small.kz", "тестовое письмо", "тест", appInfo)

	if err != nil {

		println(err)
	}
}

```