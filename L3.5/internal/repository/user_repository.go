package repository

import (
	"context"
	"database/sql"
	"errors"

	"event-booker/internal/models"

	"github.com/wb-go/wbf/dbpg"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	UpdateTelegramInfo(ctx context.Context, userID string, username string, chatID int64) error
	GetUserByTelegramUsername(ctx context.Context, username string) (*models.User, error)
}

type UserRepo struct {
	db *dbpg.DB
}

func NewUserRepository(db *dbpg.DB) UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, telegram_username, preferred_notification, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.TelegramUsername,
		user.PreferredNotification,
	)

	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return ErrEmailAlreadyExists
		}
		return err
	}

	return nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, password_hash, role, telegram_username, telegram_chat_id, 
		       preferred_notification, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.TelegramUsername,
		&user.TelegramChatID,
		&user.PreferredNotification,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, password_hash, role, telegram_username, telegram_chat_id,
		       preferred_notification, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.TelegramUsername,
		&user.TelegramChatID,
		&user.PreferredNotification,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, role = $2, telegram_username = $3, telegram_chat_id = $4, 
		    preferred_notification = $5, updated_at = NOW()
		WHERE id = $6
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Role,
		user.TelegramUsername,
		user.TelegramChatID,
		user.PreferredNotification,
		user.ID,
	)
	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return ErrEmailAlreadyExists
		}
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepo) UpdateTelegramInfo(ctx context.Context, userID string, username string, chatID int64) error {
	query := `
		UPDATE users
		SET telegram_username = $1, telegram_chat_id = $2, updated_at = NOW()
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, username, chatID, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepo) GetUserByTelegramUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, password_hash, role, telegram_username, telegram_chat_id,
		       preferred_notification, created_at, updated_at
		FROM users
		WHERE telegram_username = $1
	`

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.TelegramUsername,
		&user.TelegramChatID,
		&user.PreferredNotification,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

