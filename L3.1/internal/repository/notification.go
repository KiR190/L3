package repository

import (
	"context"

	"delayed-notifier/internal/models"
)

type NotificationRepo interface {
	Create(ctx context.Context, n *models.Notification) error
	GetByID(ctx context.Context, id int) (*models.Notification, error)
	GetActive(ctx context.Context) ([]*models.Notification, error)
	Delete(ctx context.Context, id int) error
	UpdateStatus(ctx context.Context, id int, status models.StatusType) error
	GetAll(ctx context.Context) ([]*models.Notification, error)
	UpdateRetryCount(ctx context.Context, id int, retryCount int) error
}
