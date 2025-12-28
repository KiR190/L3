package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"event-booker/internal/models"

	"github.com/wb-go/wbf/dbpg"
)

var (
	ErrEventNotFound      = errors.New("event not found")
	ErrBookingNotFound    = errors.New("booking not found")
	ErrInsufficientSeats  = errors.New("insufficient seats available")
	ErrBookingAlreadyPaid = errors.New("booking already paid")
)

type EventRepository interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEventByID(ctx context.Context, id string) (*models.Event, error)
	GetEvents(ctx context.Context) ([]models.Event, error)
}

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *models.Booking) error
	GetBookingByID(ctx context.Context, id string) (*models.Booking, error)
	GetBookingWithUser(ctx context.Context, id string) (*models.BookingWithUser, error)
	GetBookingsByEventID(ctx context.Context, eventID string) ([]models.Booking, error)
	GetBookingsByUserID(ctx context.Context, userID string) ([]models.Booking, error)
	ConfirmBooking(ctx context.Context, bookingID string) error
	CancelBooking(ctx context.Context, bookingID string) error
	GetExpiredUnpaidBookings(ctx context.Context) ([]models.Booking, error)
	GetExpiredUnpaidBookingsWithUser(ctx context.Context) ([]models.BookingWithUser, error)
	GetAvailableSeats(ctx context.Context, eventID string) (int, error)
}

type Repository struct {
	db *dbpg.DB
}

func NewRepository(db *dbpg.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateEvent(ctx context.Context, event *models.Event) error {
	query := `
		INSERT INTO events (id, name, event_date, total_seats, payment_timeout_minutes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`

	_, err := r.db.ExecContext(ctx, query,
		event.ID,
		event.Name,
		event.EventDate,
		event.TotalSeats,
		event.PaymentTimeoutMinutes,
	)

	return err
}

func (r *Repository) GetEventByID(ctx context.Context, id string) (*models.Event, error) {
	var event models.Event
	query := `
		SELECT id, name, event_date, total_seats, payment_timeout_minutes, created_at, updated_at
		FROM events
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.Name,
		&event.EventDate,
		&event.TotalSeats,
		&event.PaymentTimeoutMinutes,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrEventNotFound
	}
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *Repository) UpdateEvent(ctx context.Context, event *models.Event) error {
	query := `
		UPDATE events
		SET name = $1, event_date = $2, total_seats = $3, payment_timeout_minutes = $4, updated_at = $5
		WHERE id = $6
	`
	result, err := r.db.ExecContext(ctx, query,
		event.Name, event.EventDate, event.TotalSeats, event.PaymentTimeoutMinutes, event.UpdatedAt, event.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (r *Repository) GetEvents(ctx context.Context) ([]models.Event, error) {
	query := `
		SELECT id, name, event_date, total_seats, payment_timeout_minutes, created_at, updated_at
		FROM events
		ORDER BY event_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.EventDate,
			&event.TotalSeats,
			&event.PaymentTimeoutMinutes,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// Booking operations

func (r *Repository) CreateBooking(ctx context.Context, booking *models.Booking) error {
	// Start transaction
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Lock the event row for update
	var totalSeats int
	var timeoutMinutes int
	eventQuery := `SELECT total_seats, payment_timeout_minutes FROM events WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, eventQuery, booking.EventID).Scan(&totalSeats, &timeoutMinutes)
	if err == sql.ErrNoRows {
		return ErrEventNotFound
	}
	if err != nil {
		return err
	}

	// Count existing active bookings (unpaid and paid)
	var bookedSeats int
	countQuery := `
		SELECT COALESCE(SUM(seats_count), 0)
		FROM bookings
		WHERE event_id = $1 AND status IN ('unpaid', 'paid')
	`
	err = tx.QueryRowContext(ctx, countQuery, booking.EventID).Scan(&bookedSeats)
	if err != nil {
		return err
	}

	// Check if enough seats are available
	availableSeats := totalSeats - bookedSeats
	if availableSeats < booking.SeatsCount {
		return fmt.Errorf("%w: requested %d, available %d", ErrInsufficientSeats, booking.SeatsCount, availableSeats)
	}

	// Calculate expiration time
	booking.ExpiresAt = booking.CreatedAt.Add(time.Duration(timeoutMinutes) * time.Minute)

	// Insert the booking
	insertQuery := `
		INSERT INTO bookings (id, event_id, user_id, seats_count, status, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = tx.ExecContext(ctx, insertQuery,
		booking.ID,
		booking.EventID,
		booking.UserID,
		booking.SeatsCount,
		booking.Status,
		booking.CreatedAt,
		booking.ExpiresAt,
	)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

func (r *Repository) GetBookingByID(ctx context.Context, id string) (*models.Booking, error) {
	var booking models.Booking
	query := `
		SELECT id, event_id, user_id, seats_count, status, created_at, paid_at, expires_at
		FROM bookings
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID,
		&booking.EventID,
		&booking.UserID,
		&booking.SeatsCount,
		&booking.Status,
		&booking.CreatedAt,
		&booking.PaidAt,
		&booking.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrBookingNotFound
	}
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *Repository) GetBookingWithUser(ctx context.Context, id string) (*models.BookingWithUser, error) {
	var bwu models.BookingWithUser
	query := `
		SELECT b.id, b.event_id, b.user_id, b.seats_count, b.status, b.created_at, b.paid_at, b.expires_at,
		       u.email, u.telegram_username, u.telegram_chat_id
		FROM bookings b
		JOIN users u ON b.user_id = u.id
		WHERE b.id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&bwu.ID,
		&bwu.EventID,
		&bwu.UserID,
		&bwu.SeatsCount,
		&bwu.Status,
		&bwu.CreatedAt,
		&bwu.PaidAt,
		&bwu.ExpiresAt,
		&bwu.UserEmail,
		&bwu.TelegramUsername,
		&bwu.TelegramChatID,
	)

	if err == sql.ErrNoRows {
		return nil, ErrBookingNotFound
	}
	if err != nil {
		return nil, err
	}

	return &bwu, nil
}

func (r *Repository) GetBookingsByEventID(ctx context.Context, eventID string) ([]models.Booking, error) {
	query := `
		SELECT id, event_id, user_id, seats_count, status, created_at, paid_at, expires_at
		FROM bookings
		WHERE event_id = $1 AND status IN ('unpaid', 'paid')
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.EventID,
			&booking.UserID,
			&booking.SeatsCount,
			&booking.Status,
			&booking.CreatedAt,
			&booking.PaidAt,
			&booking.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *Repository) GetBookingsByUserID(ctx context.Context, userID string) ([]models.Booking, error) {
	query := `
		SELECT id, event_id, user_id, seats_count, status, created_at, paid_at, expires_at
		FROM bookings
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.EventID,
			&booking.UserID,
			&booking.SeatsCount,
			&booking.Status,
			&booking.CreatedAt,
			&booking.PaidAt,
			&booking.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *Repository) ConfirmBooking(ctx context.Context, bookingID string) error {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check current status
	var currentStatus string
	checkQuery := `SELECT status FROM bookings WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, checkQuery, bookingID).Scan(&currentStatus)
	if err == sql.ErrNoRows {
		return ErrBookingNotFound
	}
	if err != nil {
		return err
	}

	if currentStatus == string(models.BookingStatusPaid) {
		return ErrBookingAlreadyPaid
	}

	if currentStatus == string(models.BookingStatusCancelled) {
		return errors.New("cannot confirm cancelled booking")
	}

	// Update to paid
	updateQuery := `
		UPDATE bookings
		SET status = $1, paid_at = NOW()
		WHERE id = $2
	`

	_, err = tx.ExecContext(ctx, updateQuery, models.BookingStatusPaid, bookingID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) CancelBooking(ctx context.Context, bookingID string) error {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if booking exists and is not paid
	var currentStatus string
	checkQuery := `SELECT status FROM bookings WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, checkQuery, bookingID).Scan(&currentStatus)
	if err == sql.ErrNoRows {
		return ErrBookingNotFound
	}
	if err != nil {
		return err
	}

	// Don't cancel already paid bookings
	if currentStatus == string(models.BookingStatusPaid) {
		return errors.New("cannot cancel paid booking")
	}

	// Update status to cancelled
	updateQuery := `UPDATE bookings SET status = $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, updateQuery, models.BookingStatusCancelled, bookingID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetExpiredUnpaidBookings(ctx context.Context) ([]models.Booking, error) {
	query := `
		SELECT id, event_id, user_id, seats_count, status, created_at, paid_at, expires_at
		FROM bookings
		WHERE status = 'unpaid' AND expires_at < NOW()
		ORDER BY expires_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.EventID,
			&booking.UserID,
			&booking.SeatsCount,
			&booking.Status,
			&booking.CreatedAt,
			&booking.PaidAt,
			&booking.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *Repository) GetExpiredUnpaidBookingsWithUser(ctx context.Context) ([]models.BookingWithUser, error) {
	query := `
		SELECT b.id, b.event_id, b.user_id, b.seats_count, b.status, b.created_at, b.paid_at, b.expires_at,
		       u.email, u.telegram_username, u.telegram_chat_id
		FROM bookings b
		JOIN users u ON b.user_id = u.id
		WHERE b.status = 'unpaid' AND b.expires_at < NOW()
		ORDER BY b.expires_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.BookingWithUser
	for rows.Next() {
		var bwu models.BookingWithUser
		err := rows.Scan(
			&bwu.ID,
			&bwu.EventID,
			&bwu.UserID,
			&bwu.SeatsCount,
			&bwu.Status,
			&bwu.CreatedAt,
			&bwu.PaidAt,
			&bwu.ExpiresAt,
			&bwu.UserEmail,
			&bwu.TelegramUsername,
			&bwu.TelegramChatID,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, bwu)
	}

	return bookings, nil
}

func (r *Repository) GetAvailableSeats(ctx context.Context, eventID string) (int, error) {
	var totalSeats int
	var bookedSeats int

	// Get total seats
	eventQuery := `SELECT total_seats FROM events WHERE id = $1`
	err := r.db.QueryRowContext(ctx, eventQuery, eventID).Scan(&totalSeats)
	if err == sql.ErrNoRows {
		return 0, ErrEventNotFound
	}
	if err != nil {
		return 0, err
	}

	// Count booked seats (unpaid and paid only)
	countQuery := `
		SELECT COALESCE(SUM(seats_count), 0)
		FROM bookings
		WHERE event_id = $1 AND status IN ('unpaid', 'paid')
	`
	err = r.db.QueryRowContext(ctx, countQuery, eventID).Scan(&bookedSeats)
	if err != nil {
		return 0, err
	}

	return totalSeats - bookedSeats, nil
}

