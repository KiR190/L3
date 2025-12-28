package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"warehouse/internal/auth"
	"warehouse/internal/models"
	"warehouse/internal/repository"

	"github.com/google/uuid"
)

var ErrForbidden = errors.New("forbidden")

type Service struct {
	repo        repository.ItemRepository
	userRepo    repository.UserRepository
	authService *auth.AuthService
}

func NewService(repo repository.ItemRepository, userRepo repository.UserRepository, authService *auth.AuthService) *Service {
	return &Service{
		repo:        repo,
		userRepo:    userRepo,
		authService: authService,
	}
}

func (s *Service) CreateItem(ctx context.Context, it models.ItemCreate) (*models.Item, error) {
	item := &models.Item{
		ID:          uuid.NewString(),
		SKU:         it.SKU,
		Name:        it.Name,
		Description: it.Description,
		Quantity:    it.Quantity,
		Location:    it.Location,
	}

	err := s.repo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) GetItems(ctx context.Context, limit, offset int) ([]models.Item, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *Service) GetItem(ctx context.Context, id string) (*models.Item, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateItem(ctx context.Context, id string, upd models.ItemUpdate) (*models.Item, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if upd.SKU != nil {
		item.SKU = *upd.SKU
	}
	if upd.Name != nil {
		item.Name = *upd.Name
	}
	if upd.Description != nil {
		item.Description = *upd.Description
	}
	if upd.Quantity != nil {
		item.Quantity = *upd.Quantity
	}
	if upd.Location != nil {
		item.Location = *upd.Location
	}

	item.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) DeleteItem(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetItemHistory(ctx context.Context, itemID string, limit, offset int) ([]models.ItemHistory, error) {
	return s.repo.GetHistoryByItemID(ctx, itemID, limit, offset)
}

func canEdit(role string) bool {
	return role == "admin" || role == "manager"
}

func canDelete(role string) bool {
	return role == "admin"
}

// Auth / users

func (s *Service) Register(ctx context.Context, req models.UserRegister) (*models.AuthResponse, error) {
	passwordHash, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         models.RoleViewer,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		if err == repository.ErrEmailAlreadyExists {
			return nil, fmt.Errorf("email already registered")
		}
		return nil, err
	}

	token, err := s.authService.GenerateJWT(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		User: &models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  string(user.Role),
		},
		Token: token,
	}, nil
}

func (s *Service) Login(ctx context.Context, req models.UserLogin) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, err
	}

	err = s.authService.ComparePassword(user.PasswordHash, req.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	token, err := s.authService.GenerateJWT(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		User: &models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  string(user.Role),
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
		ID:    user.ID,
		Email: user.Email,
		Role:  string(user.Role),
	}, nil
}

func (s *Service) UpdateUserRole(ctx context.Context, userID string, role models.UserRole) error {
	return s.userRepo.UpdateUserRole(ctx, userID, role)
}
