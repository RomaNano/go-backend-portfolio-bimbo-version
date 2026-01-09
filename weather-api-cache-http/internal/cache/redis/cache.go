package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)


type Cache struct{
	client *redis.Client
}

func New(addr string) *Cache{
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &Cache{client: rdb}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (c *Cache) Set(ctx context.Context, key string, value string, ttl time.Duration) error{
	return c.client.Set(ctx, key, value, ttl).Err()
}