package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"delayed-notifier/internal/cache"
	"delayed-notifier/internal/models"
	"delayed-notifier/internal/queue"
	"delayed-notifier/internal/repository"
	"delayed-notifier/internal/sender"

	"github.com/wb-go/wbf/retry"
)

type NotificationService struct {
	repo   repository.NotificationRepo
	cache  cache.NotifCache
	queue  *queue.Queue
	sender sender.Sender
}

func NewNotificationService(repo repository.NotificationRepo, cache cache.NotifCache, queue *queue.Queue, sender sender.Sender) *NotificationService {
	return &NotificationService{
		repo:   repo,
		cache:  cache,
		queue:  queue,
		sender: sender,
	}
}

func (s *NotificationService) CreateNotification(ctx context.Context, n *models.Notification) error {
	if n.SendAt.Before(time.Now()) {
		return errors.New("send time must be in the future")
	}

	// Создаём в БД
	if err := s.repo.Create(ctx, n); err != nil {
		return err
	}

	// Сохраняем в кэш
	if err := s.cache.Set(ctx, n); err != nil {
		log.Printf("warning: failed to set notification %s in cache: %v", n.ID, err)
	}

	// Публикуем в очередь
	body, err := json.Marshal(n)
	if err != nil {
		return errors.New("failed to serialize notification")
	}

	if err := s.queue.Publish(body, n.SendAt); err != nil {
		log.Printf("failed to publish message: %v", err)
		return errors.New("failed to enqueue notification")
	}

	return nil
}

func (s *NotificationService) GetNotification(ctx context.Context, id int) (*models.Notification, error) {
	if notif, err := s.cache.Get(ctx, id); err != nil {
		return nil, err
	} else if notif != nil {
		return notif, nil
	}

	notif, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	if err := s.cache.Set(ctx, notif); err != nil {
		log.Printf("warning: failed to set notification %s in cache: %v", notif.ID, err)
	}

	return notif, nil
}

func (s *NotificationService) ListNotifications(ctx context.Context) ([]*models.Notification, error) {
	notif, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return notif, nil
}

func (s *NotificationService) CancelNotification(ctx context.Context, id int) error {
	// Удаляем из БД
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Удаляем из кэша
	if err := s.cache.Delete(ctx, id); err != nil {
		log.Printf("warning: failed to delete notification %d from cache: %v", id, err)
	}

	return nil
}

func (s *NotificationService) RestoreCacheFromDB(ctx context.Context) error {
	active, err := s.repo.GetActive(ctx)
	if err != nil {
		return err
	}

	for _, notif := range active {
		if err := s.cache.Set(ctx, notif); err != nil {
			log.Printf("не удалось добавить в кэш уведомление %s: %v", notif.ID, err)
		}
	}

	log.Printf("кэш восстановлен: добавлено %d уведомлений", len(active))
	return nil
}

func (s *NotificationService) UpdateNotificationStatus(ctx context.Context, id int, status models.StatusType) error {
	// Получаем уведомление из БД
	notif, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if notif == nil {
		return fmt.Errorf("notification %d not found", id)
	}

	// Меняем статус
	notif.Status = status
	notif.UpdatedAt = time.Now()

	// Обновляем в БД
	key, err := strconv.Atoi(notif.ID)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateStatus(ctx, key, status); err != nil {
		return err
	}

	// Обновляем в кэше
	if err := s.cache.Set(ctx, notif); err != nil {
		log.Printf("warning: failed to update notification %d in cache: %v", id, err)
	}

	return nil
}

func (s *NotificationService) IncrementRetryCount(ctx context.Context, id int) error {
	// Получаем уведомление из БД
	notif, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if notif == nil {
		return fmt.Errorf("notification %d not found", id)
	}

	// Увеличиваем счетчик ретраев
	notif.Retry++
	notif.UpdatedAt = time.Now()

	// Обновляем в базе
	if err := s.repo.UpdateRetryCount(ctx, id, notif.Retry); err != nil {
		return err
	}

	// Обновляем в кэше
	if err := s.cache.Set(ctx, notif); err != nil {
		log.Printf("warning: failed to update retry count for notification %d in cache: %v", id, err)
	}

	return nil
}

func (s *NotificationService) ProcessNotification(ctx context.Context, notif *models.Notification) error {
	id, err := strconv.Atoi(notif.ID)
	if err != nil {
		return err
	}

	// Проверяем кэш
	cachedNotif, err := s.cache.Get(ctx, id)
	if err != nil {
		log.Printf("failed to check cache for notif %d: %v", id, err)
	} else if cachedNotif == nil {
		log.Printf("notification %d not found in cache (probably canceled), skipping send", id)
		return nil
	} else if cachedNotif.Status == models.Canceled {
		log.Printf("notification %d is canceled, skipping send", id)
		return nil
	}

	// Стратегия повторов
	strategy := retry.Strategy{
		Attempts: 5, // число попыток
		Delay:    1 * time.Minute,
		Backoff:  2, // экспоненциальная задержка
	}

	// Пробуем отправить с retry
	if err := retry.Do(func() error {
		if sendErr := s.sender.Send(*notif); sendErr != nil {
			if updateErr := s.IncrementRetryCount(ctx, id); updateErr != nil {
				log.Printf("failed to increment retry count for notif %d: %v", id, updateErr)
			}
			return sendErr
		}
		return nil
	}, strategy); err != nil {
		log.Printf("notification %d failed after retries: %v", id, err)
		return s.UpdateNotificationStatus(ctx, id, models.Failed)
	}

	return s.UpdateNotificationStatus(ctx, id, models.Sent)
}
