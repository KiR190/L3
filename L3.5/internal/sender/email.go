package sender

import (
	"fmt"
	"log"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

// EmailSender handles email notifications
type EmailSender struct {
	Server   string
	Port     int
	Username string
	Password string
	From     string
}

// NewEmailSender creates a new email sender
func NewEmailSender(server string, port int, username, password, from string) *EmailSender {
	return &EmailSender{
		Server:   server,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

// Send sends an email notification
func (s *EmailSender) Send(to, subject, body string) error {
	// Configure SMTP server
	smtp := mail.NewSMTPClient()
	smtp.Host = s.Server
	smtp.Port = s.Port
	smtp.Username = s.Username
	smtp.Password = s.Password
	smtp.Encryption = mail.EncryptionSTARTTLS
	smtp.KeepAlive = false
	smtp.ConnectTimeout = 10 * time.Second
	smtp.SendTimeout = 10 * time.Second

	// Connect to SMTP server
	client, err := smtp.Connect()
	if err != nil {
		log.Printf("SMTP connection error: %v", err)
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}

	// Create email
	email := mail.NewMSG()
	email.SetFrom(s.From).
		AddTo(to).
		SetSubject(subject).
		SetBody(mail.TextPlain, body)

	// Send email
	if err := email.Send(client); err != nil {
		log.Printf("Error sending email to %s: %v", to, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("✉️  Email sent to %s: %s", to, subject)
	return nil
}

