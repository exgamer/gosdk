package structures

// EmailMessage - модель пэйлоада входящего сообщения из кафки для отправки email
type EmailMessage struct {
	ServiceName string   `json:"service_name"  validate:"required"`
	Subject     string   `json:"subject"  validate:"required"`
	Content     string   `json:"content"  validate:"required"`
	Email       string   `json:"email"  validate:"required"`
	ContentType string   `json:"content-type"`
	Attachments []string `json:"attachments"`
}
