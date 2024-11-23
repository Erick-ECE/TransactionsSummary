package email

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

// GomailService implements the EmailSender interface using the gomail library.
type GomailService struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
}

// NewGomailService creates a new instance of GomailService.
func NewGomailService(smtpHost string, smtpPort int, username, password, from string) *GomailService {
	return &GomailService{
		SMTPHost: smtpHost,
		SMTPPort: smtpPort,
		Username: username,
		Password: password,
		From:     from,
	}
}

// SendEmail sends an email using SMTP.
func (s *GomailService) SendEmail(to string, subject string, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", s.From)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body) // HTML body for styled emails

	dialer := gomail.NewDialer(s.SMTPHost, s.SMTPPort, s.Username, s.Password)
	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("could not send email: %v", err)
	}
	return nil
}
