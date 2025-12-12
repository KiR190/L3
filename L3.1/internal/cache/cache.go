package cache

import (
	"context"
	"encoding/json"
	"strconv"

	"delayed-notifier/internal/models"

	"github.com/wb-go/wbf/redis"
)

type NotifCache interface {
	Set(ctx context.Context, notif *models.Notification) error
	Get(ctx context.Context, id int) (*models.Notification, error)
	Delete(ctx context.Context, id int) error
	Close() error
}

type Cache struct {
	client *redis.Client
}

func NewCache(addr, password string, db int) *Cache {
	return &Cache{
		client: redis.New(addr, password, db),
	}
}

func (c *Cache) Get(ctx context.Context, id int) (*models.Notification, error) {
	key := strconv.Itoa(id)
	val, err := c.client.Get(ctx, key)
	if err != nil {
		if err == redis.NoMatches {
			return nil, nil // нет в кэше
		}
		return nil, err
	}

	var notif models.Notification
	if err := json.Unmarshal([]byte(val), &notif); err != nil {
		return nil, err
	}

	return &notif, nil
}

func (c *Cache) Set(ctx context.Context, notif *models.Notification) error {
	data, err := json.Marshal(notif)
	if err != nil {
		return err
	}

	key := notif.ID

	return c.client.Set(ctx, key, data)
}

func (c *Cache) Delete(ctx context.Context, id int) error {
	key := strconv.Itoa(id)
	return c.client.Del(ctx, key)
}

func (c *Cache) Close() error {
	return c.client.Close()
}
