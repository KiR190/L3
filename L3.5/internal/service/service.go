package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"event-booker/internal/auth"
	"event-booker/internal/models"
	"event-booker/internal/queue"
	"event-booker/internal/repository"
	"event-booker/internal/sender"

	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(ctx context.Context, req models.EventCreate, defaultTimeout int) (*models.Event, error)
	GetEvent(ctx context.Context, id string) (*models.EventDetail, error)
	GetEvents(ctx context.Context) ([]models.Event, error)
}

type BookingService interface {
	CreateBooking(ctx context.Context, eventID, userID string, req models.BookingRequest) (*models.Booking, error)
	ConfirmBooking(ctx context.Context, bookingID, userID string) (*models.Booking, error)
	GetBooking(ctx context.Context, bookingID string) (*models.Booking, error)
	GetUserBookings(ctx context.Context, userID string) ([]models.Booking, error)
}

type UserService interface {
	Register(ctx context.Context, req models.UserRegister) (*models.AuthResponse, error)
	Login(ctx context.Context, req models.UserLogin) (*models.AuthResponse, error)
	GetUser(ctx context.Context, userID string) (*models.UserResponse, error)
	LinkTelegram(ctx context.Context, userID string, username string) error
	UpdateTelegramChatID(ctx context.Context, username string, chatID int64) error
}

type NotificationService interface {
	NotifyBookingCancelled(ctx context.Context, booking *models.BookingWithUser, event *models.Event) error
}

type Service struct {
	repo             *repository.Repository
	userRepo         repository.UserRepository
	queue            queue.QueueInterface
	authService      *auth.AuthService
	notificationSender *sender.MultiSender
}

func NewService(repo *repository.Repository, userRepo repository.UserRepository, queue queue.QueueInterface, authService *auth.AuthService, notificationSender *sender.MultiSender) *Service {
	return &Service{
		repo:             repo,
		userRepo:         userRepo,
		queue:            queue,
		authService:      authService,
		notificationSender: notificationSender,
	}
}

// Event operations

func (s *Service) CreateEvent(ctx context.Context, req models.EventCreate, defaultTimeout int) (*models.Event, error) {
	// Use provided timeout or default
	timeoutMinutes := defaultTimeout
	if req.PaymentTimeoutMinutes != nil {
		timeoutMinutes = *req.PaymentTimeoutMinutes
	}

	event := &models.Event{
		ID:                    uuid.NewString(),
		Name:                  req.Name,
		EventDate:             req.EventDate,
		TotalSeats:            req.TotalSeats,
		PaymentTimeoutMinutes: timeoutMinutes,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	err := s.repo.CreateEvent(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	log.Printf("Created event: %s (ID: %s) with %d seats, timeout: %d minutes", event.Name, event.ID, event.TotalSeats, event.PaymentTimeoutMinutes)

	return event, nil
}

func (s *Service) GetEvent(ctx context.Context, id string) (*models.EventDetail, error) {
	// Get event
	event, err := s.repo.GetEventByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get available seats
	availableSeats, err := s.repo.GetAvailableSeats(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get active bookings
	bookings, err := s.repo.GetBookingsByEventID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.EventDetail{
		Event:          *event,
		AvailableSeats: availableSeats,
		ActiveBookings: bookings,
	}, nil
}

func (s *Service) GetEvents(ctx context.Context) ([]models.Event, error) {
	return s.repo.GetEvents(ctx)
}

func (s *Service) UpdateEvent(ctx context.Context, eventID string, req models.EventUpdate) (*models.Event, error) {
	// Get existing event
	event, err := s.repo.GetEventByID(ctx, eventID)
	if err != nil {
		if err == repository.ErrEventNotFound {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		event.Name = *req.Name
	}
	if req.EventDate != nil {
		event.EventDate = *req.EventDate
	}
	if req.TotalSeats != nil {
		event.TotalSeats = *req.TotalSeats
	}
	if req.PaymentTimeoutMinutes != nil {
		event.PaymentTimeoutMinutes = *req.PaymentTimeoutMinutes
	}

	event.UpdatedAt = time.Now()

	err = s.repo.UpdateEvent(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	log.Printf("Event %s updated", eventID)
	return event, nil
}

// Booking operations

func (s *Service) CreateBooking(ctx context.Context, eventID, userID string, req models.BookingRequest) (*models.Booking, error) {
	// Validate seats count
	if req.SeatsCount <= 0 {
		return nil, fmt.Errorf("seats count must be greater than 0")
	}

	booking := &models.Booking{
		ID:         uuid.NewString(),
		EventID:    eventID,
		UserID:     userID,
		SeatsCount: req.SeatsCount,
		Status:     models.BookingStatusUnpaid,
		CreatedAt:  time.Now(),
	}

	// Create booking (this includes transaction-safe seat availability check)
	err := s.repo.CreateBooking(ctx, booking)
	if err != nil {
		return nil, err
	}

	log.Printf("Created booking: %s for event %s (%d seats, expires at %s)", booking.ID, eventID, booking.SeatsCount, booking.ExpiresAt.Format(time.RFC3339))

	// Publish expiration task to queue
	expirationTask := models.ExpirationTask{
		BookingID: booking.ID,
		ExpiresAt: booking.ExpiresAt,
	}

	err = s.queue.PublishExpiration(ctx, expirationTask)
	if err != nil {
		log.Printf("Warning: failed to publish expiration task for booking %s: %v", booking.ID, err)
		// Don't fail the booking creation if queue publish fails
	}

	return booking, nil
}

func (s *Service) ConfirmBooking(ctx context.Context, bookingID, userID string) (*models.Booking, error) {
	// Get booking first to verify ownership
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	// Verify the booking belongs to the user
	if booking.UserID != userID {
		return nil, fmt.Errorf("booking does not belong to user")
	}

	// Confirm the booking
	err = s.repo.ConfirmBooking(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	// Get updated booking
	booking, err = s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	log.Printf("Confirmed payment for booking: %s (event: %s, user: %s)", bookingID, booking.EventID, userID)

	return booking, nil
}

func (s *Service) GetBooking(ctx context.Context, bookingID string) (*models.Booking, error) {
	return s.repo.GetBookingByID(ctx, bookingID)
}

func (s *Service) GetUserBookings(ctx context.Context, userID string) ([]models.Booking, error) {
	return s.repo.GetBookingsByUserID(ctx, userID)
}

// User operations

func (s *Service) Register(ctx context.Context, req models.UserRegister) (*models.AuthResponse, error) {
	// Hash password
	passwordHash, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ID:                    uuid.NewString(),
		Email:                 req.Email,
		PasswordHash:          passwordHash,
		Role:                  models.RoleUser, // Default role
		TelegramUsername:      req.TelegramUsername,
		PreferredNotification: "email",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		if err == repository.ErrEmailAlreadyExists {
			return nil, fmt.Errorf("email already registered")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.authService.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	log.Printf("User registered: %s (%s)", user.Email, user.ID)

	return &models.AuthResponse{
		User: &models.UserResponse{
			ID:                   user.ID,
			Email:                user.Email,
			Role:                 string(user.Role),
			TelegramUsername:     user.TelegramUsername,
			TelegramRegistered:   false,
			PreferredNotification: user.PreferredNotification,
		},
		Token: token,
	}, nil
}

func (s *Service) Login(ctx context.Context, req models.UserLogin) (*models.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Compare password
	err = s.authService.ComparePassword(user.PasswordHash, req.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := s.authService.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	log.Printf("User logged in: %s (%s)", user.Email, user.ID)

	return &models.AuthResponse{
		User: &models.UserResponse{
			ID:                   user.ID,
			Email:                user.Email,
			Role:                 string(user.Role),
			TelegramUsername:     user.TelegramUsername,
			TelegramRegistered:   user.TelegramChatID != nil && *user.TelegramChatID != 0,
			PreferredNotification: user.PreferredNotification,
		},
		Token: token,
	}, nil
}

func (s *Service) GetUser(ctx context.Context, userID string) (*models.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:                   user.ID,
		Email:                user.Email,
		Role:                 string(user.Role),
		TelegramUsername:     user.TelegramUsername,
		TelegramRegistered:   user.TelegramChatID != nil && *user.TelegramChatID != 0,
		PreferredNotification: user.PreferredNotification,
	}, nil
}

func (s *Service) LinkTelegram(ctx context.Context, userID string, username string) error {
	// Update user's telegram username
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.TelegramUsername = &username
	user.PreferredNotification = "telegram"
	user.UpdatedAt = time.Now()

	err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Note: The actual chat ID will be set when user sends /register to the bot
	log.Printf("User %s linked Telegram username: @%s and set preferred notification to telegram", userID, username)

	return nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID string, req models.UpdateProfile) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if req.Email != nil && *req.Email != "" {
		user.Email = *req.Email
	}
	if req.PreferredNotification != nil {
		user.PreferredNotification = *req.PreferredNotification
	}
	if req.TelegramUsername != nil {
		user.TelegramUsername = req.TelegramUsername
	}

	user.UpdatedAt = time.Now()

	err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		if err == repository.ErrEmailAlreadyExists {
			return nil, fmt.Errorf("email already in use")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	log.Printf("User %s profile updated", userID)
	return user, nil
}

func (s *Service) UpdateTelegramChatID(ctx context.Context, username string, chatID int64) error {
	// Find user by telegram username
	user, err := s.userRepo.GetUserByTelegramUsername(ctx, username)
	if err != nil {
		if err == repository.ErrUserNotFound {
			// User not found, this is okay - they might register in bot before web
			log.Printf("Telegram user @%s registered in bot but not in system yet", username)
			return nil
		}
		return err
	}

	// Update telegram chat ID
	err = s.userRepo.UpdateTelegramInfo(ctx, user.ID, username, chatID)
	if err != nil {
		return err
	}

	log.Printf("Updated Telegram chat ID for user %s (@%s): %d", user.Email, username, chatID)
	return nil
}

// Notification operations

func (s *Service) NotifyBookingCancelled(ctx context.Context, booking *models.BookingWithUser, event *models.Event) error {
	subject := "⚠️ Booking Cancelled"
	body := fmt.Sprintf(`Your booking has been automatically cancelled.

Event: %s
Date: %s
Seats: %d
Reason: Payment timeout

The seats have been released and are now available for other customers.`,
		event.Name,
		event.EventDate.Format("January 2, 2006 at 3:04 PM"),
		booking.SeatsCount,
	)

	err := s.notificationSender.SendWithPriority(
		booking.TelegramUsername,
		booking.TelegramUsername,
		booking.TelegramChatID,
		booking.UserEmail,
		subject,
		body,
	)

	if err != nil {
		log.Printf("Failed to send cancellation notification to user %s: %v", booking.UserEmail, err)
		return err
	}

	log.Printf("Sent cancellation notification to user %s for booking %s", booking.UserEmail, booking.ID)
	return nil
}

