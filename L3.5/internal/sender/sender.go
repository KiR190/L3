package sender

import (
	"fmt"
	"log"
)

// Sender interface for notification senders
type Sender interface {
	Send(to, subject, body string) error
}

// MultiSender sends notifications with priority fallback (Telegram → Email)
type MultiSender struct {
	emailSender    *EmailSender
	telegramSender *TelegramSender
}

// NewMultiSender creates a new multi-sender with priority fallback
func NewMultiSender(emailSender *EmailSender, telegramSender *TelegramSender) *MultiSender {
	return &MultiSender{
		emailSender:    emailSender,
		telegramSender: telegramSender,
	}
}

// SendWithPriority tries Telegram first, falls back to Email
// telegramUsername: username without @ prefix (can be empty)
// email: email address for fallback
func (m *MultiSender) SendWithPriority(telegramUsername, telegramChatID *string, chatID *int64, email, subject, body string) error {
	var telegramErr error
	
	// Try Telegram first if chat ID is available
	if chatID != nil && *chatID != 0 {
		telegramErr = m.telegramSender.SendToChat(*chatID, subject, body)
		if telegramErr == nil {
			return nil // Success via Telegram
		}
		log.Printf("Telegram delivery failed (chatID: %d): %v, falling back to email", *chatID, telegramErr)
	} else if telegramUsername != nil && *telegramUsername != "" {
		// Try with username if no chat ID
		telegramErr = m.telegramSender.Send(*telegramUsername, subject, body)
		if telegramErr == nil {
			return nil // Success via Telegram
		}
		log.Printf("Telegram delivery failed (@%s): %v, falling back to email", *telegramUsername, telegramErr)
	}

	// Fallback to email
	if m.emailSender != nil && m.emailSender.Server != "" {
		emailErr := m.emailSender.Send(email, subject, body)
		if emailErr == nil {
			return nil // Success via Email
		}
		
		// Both failed
		return fmt.Errorf("all delivery methods failed - Telegram: %v, Email: %v", telegramErr, emailErr)
	}

	// Only Telegram was attempted and it failed
	if telegramErr != nil {
		return fmt.Errorf("telegram delivery failed and no email configured: %v", telegramErr)
	}

	return fmt.Errorf("no notification channels available")
}

// Send implements the Sender interface (sends only to email)
func (m *MultiSender) Send(to, subject, body string) error {
	if m.emailSender != nil && m.emailSender.Server != "" {
		return m.emailSender.Send(to, subject, body)
	}
	return fmt.Errorf("email sender not configured")
}

