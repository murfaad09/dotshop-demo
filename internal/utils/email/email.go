package email

import (
	"fmt"
	"net/smtp"
)

// SMTPConfig holds the configuration for the SMTP server
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// EmailSender is responsible for sending emails
type EmailSender struct {
	config SMTPConfig
}

// NewEmailSender creates a new EmailSender
func NewEmailSender(config SMTPConfig) *EmailSender {
	return &EmailSender{config: config}
}

// SendEmail sends an email
func (es *EmailSender) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", es.config.Username, es.config.Password, es.config.Host)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))
	addr := fmt.Sprintf("%s:%s", es.config.Host, es.config.Port)

	err := smtp.SendMail(addr, auth, es.config.From, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}
