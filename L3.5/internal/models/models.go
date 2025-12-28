package models

import "time"

// BookingStatus represents the state of a booking
type BookingStatus string

const (
	BookingStatusUnpaid    BookingStatus = "unpaid"
	BookingStatusPaid      BookingStatus = "paid"
	BookingStatusCancelled BookingStatus = "cancelled"
)

// Event represents an event that can be booked
type Event struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	EventDate             time.Time `json:"event_date"`
	TotalSeats            int       `json:"total_seats"`
	PaymentTimeoutMinutes int       `json:"payment_timeout_minutes"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// Booking represents a seat reservation
type Booking struct {
	ID         string        `json:"id"`
	EventID    string        `json:"event_id"`
	UserID     string        `json:"user_id"`
	SeatsCount int           `json:"seats_count"`
	Status     BookingStatus `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	PaidAt     *time.Time    `json:"paid_at,omitempty"`
	ExpiresAt  time.Time     `json:"expires_at"`
}

// UserRole represents user role type
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// User represents a registered user
type User struct {
	ID                    string     `json:"id"`
	Email                 string     `json:"email"`
	PasswordHash          string     `json:"-"` // Never expose in JSON
	Role                  UserRole   `json:"role"`
	TelegramUsername      *string    `json:"telegram_username,omitempty"`
	TelegramChatID        *int64     `json:"-"` // Internal use only
	PreferredNotification string     `json:"preferred_notification"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// BookingWithUser represents a booking with associated user information
type BookingWithUser struct {
	Booking
	UserEmail         string  `json:"user_email"`
	TelegramUsername  *string `json:"telegram_username,omitempty"`
	TelegramChatID    *int64  `json:"-"`
}

// EventCreate is the DTO for creating a new event
type EventCreate struct {
	Name                  string    `json:"name" binding:"required,min=1,max=200"`
	EventDate             time.Time `json:"event_date" binding:"required"`
	TotalSeats            int       `json:"total_seats" binding:"required,gt=0,lte=100000"`
	PaymentTimeoutMinutes *int      `json:"payment_timeout_minutes,omitempty" binding:"omitempty,gt=0,lte=1440"`
}

// BookingRequest is the DTO for creating a new booking
type BookingRequest struct {
	SeatsCount int `json:"seats_count" binding:"required,gt=0,lte=100"`
}

// EventDetail includes event info with real-time availability
type EventDetail struct {
	Event           Event     `json:"event"`
	AvailableSeats  int       `json:"available_seats"`
	ActiveBookings  []Booking `json:"active_bookings"`
}

// BookingConfirmRequest is the DTO for confirming a booking payment
type BookingConfirmRequest struct {
	BookingID string `json:"booking_id" validate:"required"`
}

// ExpirationTask represents a message for the worker queue
type ExpirationTask struct {
	BookingID string    `json:"booking_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// UserRegister is the DTO for user registration
type UserRegister struct {
	Email            string  `json:"email" binding:"required,email"`
	Password         string  `json:"password" binding:"required,min=8"`
	TelegramUsername *string `json:"telegram_username,omitempty" binding:"omitempty,min=1,max=32"`
}

// UserLogin is the DTO for user login
type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateProfile is the DTO for updating user profile
type UpdateProfile struct {
	Email                 *string `json:"email,omitempty" binding:"omitempty,email"`
	PreferredNotification *string `json:"preferred_notification,omitempty" binding:"omitempty,oneof=email telegram"`
	TelegramUsername      *string `json:"telegram_username,omitempty" binding:"omitempty,min=1,max=32"`
}

// AuthResponse is the response after successful auth
type AuthResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

// UserResponse is the public user info
type UserResponse struct {
	ID                   string  `json:"id"`
	Email                string  `json:"email"`
	Role                 string  `json:"role"`
	TelegramUsername     *string `json:"telegram_username,omitempty"`
	TelegramRegistered   bool    `json:"telegram_registered"`
	PreferredNotification string `json:"preferred_notification"`
}

// EventUpdate is the DTO for updating an event
type EventUpdate struct {
	Name                  *string    `json:"name,omitempty" binding:"omitempty,min=1,max=200"`
	EventDate             *time.Time `json:"event_date,omitempty"`
	TotalSeats            *int       `json:"total_seats,omitempty" binding:"omitempty,gt=0,lte=100000"`
	PaymentTimeoutMinutes *int       `json:"payment_timeout_minutes,omitempty" binding:"omitempty,gt=0,lte=1440"`
}

// LinkTelegramRequest is the DTO for linking Telegram account
type LinkTelegramRequest struct {
	TelegramUsername string `json:"telegram_username" binding:"required,min=1,max=32"`
}

