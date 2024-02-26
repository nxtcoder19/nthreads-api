package redis

import (
	"context"
	r "github.com/redis/go-redis/v9"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value []byte, exp *time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
}

type CacheImpl struct {
	client *r.Client
}

func (c CacheImpl) Set(ctx context.Context, key string, value []byte, exp *time.Duration) error {
	if exp == nil {
		return c.client.Set(ctx, key, value, -1).Err()
	}
	return c.client.Set(ctx, key, value, *exp).Err()
}

func (c CacheImpl) Get(ctx context.Context, key string) ([]byte, error) {
	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}

func (c CacheImpl) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func NewRedisCache(addr string) Cache {
	client := r.NewClient(&r.Options{
		Addr: addr,
	})
	return CacheImpl{client: client}
}
