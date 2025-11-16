package models

import "time"

type ShortURL struct {
	ID        int
	ShortCode string
	Original  string
	CreatedAt time.Time
}

type ClickEvent struct {
	ID        int
	ShortID   int
	UserAgent string
	Timestamp time.Time
}

