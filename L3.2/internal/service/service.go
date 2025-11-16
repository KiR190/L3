package service

import (
	"context"
	"fmt"
	"log"

	"shortener/internal/cache"
	"shortener/internal/generator"
	"shortener/internal/models"
	. "shortener/internal/repository"

	"github.com/lib/pq"
)

type ShortenerService struct {
	shortRepo ShortURLRepository
	analytics AnalyticsRepository
	cache     cache.Cache
	generator generator.ShortCodeGenerator
}

func NewShortenerService(s ShortURLRepository, a AnalyticsRepository, c cache.Cache, g generator.ShortCodeGenerator) *ShortenerService {
	return &ShortenerService{s, a, c, g}
}

func (s *ShortenerService) Create(ctx context.Context, original, customCode string) (*models.ShortURL, error) {
	if customCode != "" {
		url := &models.ShortURL{
			ShortCode: customCode,
			Original:  original,
		}

		err := s.shortRepo.Save(ctx, url)
		if err != nil {
			if isUniqueViolation(err) {
				return nil, fmt.Errorf("short code '%s' is already taken", customCode)
			}
			return nil, err
		}

		return url, nil
	}

	for i := 0; i < 3; i++ {
		code := s.generator.Generate()

		url := &models.ShortURL{
			ShortCode: code,
			Original:  original,
		}

		err := s.shortRepo.Save(ctx, url)
		if err == nil {
			return url, nil
		}

		// если конфликт — пробуем ещё раз
		if isUniqueViolation(err) {
			continue
		}
		return nil, err
	}
	return nil, fmt.Errorf("failed to create short url after retries")
}

func isUniqueViolation(err error) bool {
	pgErr, ok := err.(*pq.Error)
	return ok && pgErr.Code == "23505"
}

func (s *ShortenerService) Resolve(ctx context.Context, code, ua string) (string, error) {
	if cached, err := s.cache.Get(ctx, code); err != nil {
		return "", err
	} else if cached != nil {
		s.recordClick(cached.ID, ua)
		return cached.Original, nil
	}

	url, err := s.shortRepo.FindByID(ctx, code)
	if err != nil {
		return "", err
	}

	_ = s.cache.Set(ctx, url)

	s.recordClick(url.ID, ua)
	return url.Original, nil
}

func (s *ShortenerService) recordClick(urlID int, ua string) {
	click := models.ClickEvent{
		ShortID:   urlID,
		UserAgent: ua,
	}

	go s.analytics.Save(context.Background(), &click)
}

func (s *ShortenerService) GetAnalytics(ctx context.Context, code string) ([]*models.ClickEvent, error) {
	return s.analytics.GetStats(ctx, code)
}

func (s *ShortenerService) ListLatest(ctx context.Context) ([]models.ShortURL, error) {
	return s.shortRepo.ListLatest(ctx, 100)
}

func (s *ShortenerService) RestoreCacheFromDB(ctx context.Context) error {
	urls, err := s.shortRepo.FindTopPopular(ctx, 100)
	if err != nil {
		return err
	}

	for _, url := range urls {
		if err := s.cache.Set(ctx, &url); err != nil {
			log.Printf("не удалось добавить в кэш %d: %v", url.ID, err)
		}
	}

	log.Printf("кэш восстановлен: добавлено %d ссылок", len(urls))
	return nil
}
