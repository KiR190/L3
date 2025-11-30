package service

import (
	"context"
	"errors"
	"time"

	"sales-tracker/internal/models"
	"sales-tracker/internal/repository"

	"github.com/google/uuid"
)

type ItemService interface {
	CreateItem(ctx context.Context, item models.ItemCreate) (*models.Item, error)
	GetItems(ctx context.Context) ([]models.Item, error)
	UpdateItem(ctx context.Context, id string, upd models.ItemUpdate) (*models.Item, error)
	DeleteItem(ctx context.Context, id string) error

	GetAnalytics(ctx context.Context, from, to string) (models.Analytics, error)
}

type Service struct {
	repo repository.ItemRepository
}

func NewService(repo repository.ItemRepository) ItemService {
	return &Service{repo: repo}
}

func (s *Service) CreateItem(ctx context.Context, payload models.ItemCreate) (*models.Item, error) {
	item := &models.Item{
		ID:          uuid.NewString(),
		Type:        models.ItemType(payload.Type),
		Amount:      int64(payload.Amount * 100),
		Currency:    payload.Currency,
		CategoryID:  payload.CategoryID,
		Description: payload.Description,
		OccurredAt:  payload.Date,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) GetItems(ctx context.Context) ([]models.Item, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateItem(ctx context.Context, id string, upd models.ItemUpdate) (*models.Item, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("item not found")
	}

	if upd.Type != nil {
		item.Type = models.ItemType(*upd.Type)
	}

	if upd.Amount != nil {
		item.Amount = int64(*upd.Amount * 100)
	}

	if upd.CategoryID != nil {
		item.CategoryID = upd.CategoryID
	}

	if upd.Description != nil {
		item.Description = upd.Description
	}

	if upd.Currency != nil {
		item.Currency = *upd.Currency
	}

	if upd.Date != nil {
		item.OccurredAt = *upd.Date
	}

	item.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) DeleteItem(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetAnalytics(ctx context.Context, from, to string) (models.Analytics, error) {
	return s.repo.GetAnalytics(ctx, from, to)
}
