package models

import "time"

type StatusType int

const (
	Scheduled StatusType = iota
	Sent
	Failed
	Canceled
)

type ChannelType int

const (
	Native ChannelType = iota
	Email
	Telegram
)

type Notification struct {
	ID        string      `json:"id" validate:"required"`
	UserID    string      `json:"user_id" validate:"required"`
	Channel   ChannelType `json:"channel" validate:"required"`
	Recipient string      `json:"recipient" validate:"required"`
	Message   string      `json:"message" validate:"required"`
	SendAt    time.Time   `json:"sent_at" validate:"required"`
	Status    StatusType  `json:"status" validate:"required"`
	Retry     int         `json:"retry" validate:"required"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
