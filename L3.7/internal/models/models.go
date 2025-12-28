package models

import (
	"encoding/json"
	"time"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleViewer  UserRole = "viewer"
)

type Item struct {
	ID          string    `db:"id" json:"id"`
	SKU         string    `db:"sku" json:"sku"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description,omitempty" json:"description,omitempty"`
	Quantity    int       `db:"quantity" json:"quantity"`
	Location    string    `db:"location,omitempty" json:"location,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type ItemCreate struct {
	SKU         string `json:"sku" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity" validate:"gte=0"`
	Location    string `json:"location"`
}

type ItemUpdate struct {
	SKU         *string `json:"sku"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Quantity    *int    `json:"quantity"`
	Location    *string `json:"location"`
}

type ItemHistory struct {
	ID        string          `db:"id" json:"id"`
	ItemID    string          `db:"item_id" json:"item_id"`
	Action    string          `db:"action" json:"action"`
	OldData   json.RawMessage `db:"old_data" json:"old_data,omitempty"`
	NewData   json.RawMessage `db:"new_data" json:"new_data,omitempty"`
	UserID    *string         `db:"user_id" json:"user_id,omitempty"`
	Username  *string         `db:"username" json:"username,omitempty"`
	Role      *string         `db:"role" json:"role,omitempty"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
}

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRegister struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AuthResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin manager viewer"`
}
