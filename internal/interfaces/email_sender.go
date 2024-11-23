package interfaces

// EmailSender defines the interface for sending emails.
type EmailSender interface {
	SendEmail(to string, subject string, body string) error
}
