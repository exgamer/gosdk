package producer

import (
	"encoding/json"
	"errors"
	"github.com/exgamer/gosdk/pkg/config"
	"github.com/exgamer/gosdk/pkg/constants"
	"github.com/exgamer/gosdk/pkg/kafkahelper/messenger/structures"
	"github.com/exgamer/gosdk/pkg/logger"
	"github.com/exgamer/gosdk/pkg/validation"
)

func NewMessageSender(producer KafkaProducer) *MessageSender {
	return &MessageSender{
		producer: producer,
	}
}

// MessageSender - отсылка сообщений через кафку
type MessageSender struct {
	producer KafkaProducer
}

func (s *MessageSender) SendSms(phone string, text string, appInfo *config.AppInfo) error {
	normalizedPhone, err := validation.NormalizePhoneNumber(phone)

	if err != nil {
		message := "Sms message send error: invalid phone number:" + phone

		if appInfo != nil {
			logger.FormattedError(appInfo.ServiceName, appInfo.RequestMethod, appInfo.RequestUrl, 0, appInfo.RequestId, message)
		} else {
			println(message)
		}

		return err
	}

	smsMessage := structures.SmsMessage{
		Phone:       normalizedPhone,
		Text:        text,
		ServiceName: appInfo.ServiceName,
	}

	topic := "messenger-service.command.sms" // хард код потому что по идее никогда не изменится
	jsonValue, _ := json.Marshal(smsMessage)
	sendErr := s.producer.SendMessage(topic, string(jsonValue))

	if sendErr != nil {
		return sendErr
	}

	message := "Sms message send: phone:" + phone + "; text:" + text

	if appInfo != nil {
		logger.FormattedInfo(appInfo.ServiceName, appInfo.RequestMethod, appInfo.RequestUrl, 0, appInfo.RequestId, message)
	} else {
		println(message)
	}

	return nil
}

func (s *MessageSender) SendEmailWithHtml(email string, subject string, content string, appInfo *config.AppInfo) error {
	return s.SendEmail(email, subject, content, nil, constants.ContentTypeHtml, appInfo)
}

func (s *MessageSender) SendEmailWithPlainText(email string, subject string, content string, appInfo *config.AppInfo) error {
	return s.SendEmail(email, subject, content, nil, constants.ContentTypeText, appInfo)
}

func (s *MessageSender) SendEmail(email string, subject string, content string, attachments []string, contentType string, appInfo *config.AppInfo) error {
	emailMessage := structures.EmailMessage{
		Email:       email,
		Subject:     subject,
		Content:     content,
		ContentType: contentType,
		ServiceName: appInfo.ServiceName,
		Attachments: attachments,
	}

	if !validation.CheckValidEmail(emailMessage.Email) {
		message := "Email message send error: invalid email:" + emailMessage.Email

		if appInfo != nil {
			logger.FormattedError(appInfo.ServiceName, appInfo.RequestMethod, appInfo.RequestUrl, 0, appInfo.RequestId, message)
		} else {
			println(message)
		}

		return errors.New("invalid email")
	}

	topic := "messenger-service.command.email" // хард код потому что по идее никогда не изменится
	jsonValue, _ := json.Marshal(emailMessage)
	sendErr := s.producer.SendMessage(topic, string(jsonValue))

	if sendErr != nil {
		return sendErr
	}

	message := "Email send: email:" + emailMessage.Email + "; content:" + emailMessage.Content

	if appInfo != nil {
		logger.FormattedInfo(appInfo.ServiceName, appInfo.RequestMethod, appInfo.RequestUrl, 0, appInfo.RequestId, message)
	} else {
		println(message)
	}

	return nil
}
