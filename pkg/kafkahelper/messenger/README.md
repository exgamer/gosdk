Пример использования функции для отправки смс через кафку

```go
package app

"github.com/exgamer/gosdk/pkg/kafkahelper/messenger"
kafkaLib "github.com/confluentinc/confluent-kafkahelper-go/v2/kafkahelper"

func main() {
	kafkaProducer, _ := kafkaLib.NewProducer(&kafkaLib.ConfigMap{
		"bootstrap.servers": strings.Join([]string{app.kafkaConfig.Host}, ","),
	})

	messageSender := messenger.NewMessageSender(kafkaProducer)

	defer kafkaProducer.Close()

	err := messageSender.SendSms("+7 (778)-300-00-44", "test message")

	if err != nil {

		println(err)
	}
}

```


Пример использования функции для отправки email через кафку

```go
package app

"github.com/exgamer/gosdk/pkg/kafkahelper/messenger"
kafkaLib "github.com/confluentinc/confluent-kafkahelper-go/v2/kafkahelper"

func main() {
	kafkaProducer, _ := kafkaLib.NewProducer(&kafkaLib.ConfigMap{
		"bootstrap.servers": strings.Join([]string{app.kafkaConfig.Host}, ","),
	})

	messageSender := messenger.NewMessageSender(kafkaProducer)

	defer kafkaProducer.Close()

	err := handler.messageSender.SendEmailWithPlainText("example@small.kz", "тестовое письмо", "тест", appInfo)

	if err != nil {

		println(err)
	}
}

```