package sender

import (
	"log"

	"delayed-notifier/internal/models"
)

type NativeSender struct{}

func (n *NativeSender) Send(notification models.Notification) error {
	log.Printf("📨 Получено сообщение: %+v\n", notification)
	return nil
}
