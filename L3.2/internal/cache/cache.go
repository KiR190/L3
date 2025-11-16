package cache

import (
	"context"
	"encoding/json"
	"time"

	"shortener/internal/models"

	"github.com/wb-go/wbf/redis"
)

type Cache interface {
	Set(ctx context.Context, url *models.ShortURL) error
	Get(ctx context.Context, id string) (*models.ShortURL, error)
	Delete(ctx context.Context, id string) error
}

type URLCache struct {
	client *redis.Client
}

func NewCache(addr, password string, db int) *URLCache {
	return &URLCache{
		client: redis.New(addr, password, db),
	}
}

func (c *URLCache) Set(ctx context.Context, url *models.ShortURL) error {
	data, err := json.Marshal(url)
	if err != nil {
		return err
	}

	key := url.ShortCode
	return c.client.SetWithExpiration(ctx, key, data, time.Hour)
}

func (c *URLCache) Get(ctx context.Context, code string) (*models.ShortURL, error) {
	key := code
	val, err := c.client.Get(ctx, key)
	if err != nil {
		if err == redis.NoMatches {
			return nil, nil // нет в кэше
		}
		return nil, err
	}

	var url models.ShortURL
	if err := json.Unmarshal([]byte(val), &url); err != nil {
		return nil, err
	}

	return &url, nil
}

func (c *URLCache) Delete(ctx context.Context, code string) error {
	return c.client.Del(ctx, code)
}
