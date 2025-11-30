package models

import "time"

type ItemType string

const (
	ItemTypeIncome  ItemType = "income"
	ItemTypeExpense ItemType = "expense"
)

type Item struct {
	ID          string    `json:"id"`     // uuid
	Type        ItemType  `json:"type"`   // income/expense
	Amount      int64     `json:"amount"` // в копейках
	Currency    string    `json:"currency"`
	CategoryID  *string   `json:"category_id"`
	Description *string   `json:"description"`
	OccurredAt  time.Time `json:"occurred_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Создание
type ItemCreate struct {
	Type        string    `json:"type" validate:"required,oneof=income expense"`
	CategoryID  *string   `json:"category_id" validate:"omitempty"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	Currency    string    `json:"currency,omitempty" validate:"omitempty,len=3"`
	Description *string   `json:"description,omitempty" validate:"omitempty,max=500"`
	Date        time.Time `json:"date" validate:"required"`
}

// Обновление
type ItemUpdate struct {
	Type        *string    `json:"type,omitempty" validate:"omitempty,oneof=income expense"`
	CategoryID  *string    `json:"category_id,omitempty" validate:"omitempty"`
	Amount      *float64   `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Currency    *string    `json:"currency,omitempty" validate:"omitempty,len=3"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=500"`
	Date        *time.Time `json:"date,omitempty"`
}

type Analytics struct {
	Sum    int64
	Avg    float64
	Count  int64
	Median float64
	P90    float64
}
