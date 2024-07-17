Пример использования консьюмера для кафки

```go
package app

func (app *App) RunConsumers() {
	brokerList := []string{app.kafkaConfig.Host}
	handlers := map[string]kafka.IKafkaHandler{
		"messenger-service.command.sms": kafka2.NewKafkaSmsSaramaHandler(app.servicesFactory.SmsTrafficService),
	}

	kafka.RunConsumer(app.baseConfig, brokerList, handlers)
}


//обработчик сообщения кафки
func NewKafkaSmsSaramaHandler(
	smsSenderService contracts.ICanSendSmsService,
) *KafkaSmsSaramaHandler {
	return &KafkaSmsSaramaHandler{
		smsSenderService: smsSenderService,
	}
}

// KafkaSmsSaramaHandler - обработчик для отправки смс сообщений
type KafkaSmsSaramaHandler struct {
	smsSenderService contracts.ICanSendSmsService
}

func (h KafkaSmsSaramaHandler) Handle(message *sarama.ConsumerMessage) error {
	smsMessage := &structures2.SmsMessage{}
	err := json.Unmarshal(message.Value, smsMessage)

	if err != nil {
		return err
	}

	spew.Dump(smsMessage)
	v := validate.Struct(smsMessage)

	if !v.Validate() {
		return errors.New("payload validation error")
	}

	normalizedPhone, err := validation.NormalizePhoneNumber(smsMessage.Phone)

	if err != nil {
		return err
	}

	smsMessage.Phone = normalizedPhone
	//smsErr := h.smsSenderService.SendSms(smsMessage)
	//
	//if smsErr != nil {
	//	return smsErr
	//}

	return nil
}
```