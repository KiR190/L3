package repository

import (
	"context"

	"delayed-notifier/internal/models"

	"github.com/wb-go/wbf/dbpg"
)

type PostgresNotificationRepo struct {
	DB *dbpg.DB
}

func NewPostgresNotificationRepo(db *dbpg.DB) *PostgresNotificationRepo {
	return &PostgresNotificationRepo{
		DB: db,
	}
}

func (r *PostgresNotificationRepo) Create(ctx context.Context, n *models.Notification) error {
	query := `
		INSERT INTO notifications(user_id, channel, recipient, message, send_at)
		VALUES($1,$2,$3,$4,$5)
		RETURNING id, status, retry_count, created_at, updated_at
	`
	err := r.DB.QueryRowContext(ctx, query,
		n.UserID, n.Channel, n.Recipient, n.Message, n.SendAt,
	).Scan(&n.ID, &n.Status, &n.Retry, &n.CreatedAt, &n.UpdatedAt)
	return err
}

func (r *PostgresNotificationRepo) GetByID(ctx context.Context, id int) (*models.Notification, error) {
	query := `SELECT id, user_id, channel, recipient, message, send_at, status, retry_count, created_at, updated_at FROM notifications WHERE id=$1`
	row := r.DB.QueryRowContext(ctx, query, id)

	var n models.Notification
	err := row.Scan(&n.ID, &n.UserID, &n.Channel, &n.Recipient, &n.Message, &n.SendAt, &n.Status, &n.Retry, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *PostgresNotificationRepo) GetActive(ctx context.Context) ([]*models.Notification, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, user_id, channel, recipient, message, send_at, status, retry_count, created_at, updated_at
		FROM notifications
		WHERE status = $1
		AND send_at >= NOW()
	`, models.Scheduled)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Channel, &n.Recipient, &n.Message, &n.SendAt, &n.Status, &n.Retry, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, &n)
	}

	return result, rows.Err()
}

func (r *PostgresNotificationRepo) GetAll(ctx context.Context) ([]*models.Notification, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, user_id, channel, recipient, message, send_at, status, retry_count, created_at, updated_at
		FROM notifications
		ORDER BY created_at DESC
		LIMIT 100;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Channel, &n.Recipient, &n.Message, &n.SendAt, &n.Status, &n.Retry, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, &n)
	}

	return result, rows.Err()
}

func (r *PostgresNotificationRepo) Delete(ctx context.Context, id int) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE notifications
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`, models.Canceled, id)
	return err
}

func (r *PostgresNotificationRepo) UpdateStatus(ctx context.Context, id int, status models.StatusType) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE notifications
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`, status, id)
	return err
}

func (r *PostgresNotificationRepo) UpdateRetryCount(ctx context.Context, id int, retryCount int) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE notifications
		SET retry_count = $1, updated_at = NOW()
		WHERE id = $2
	`, retryCount, id)
	return err
}
