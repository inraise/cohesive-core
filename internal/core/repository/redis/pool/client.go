package core_pool_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client interface {
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) (bool, error)
	Close() error
}

type client struct {
	rdb *redis.Client
}

func NewClient(cfg Config) (Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &client{rdb: rdb}, nil
}

func (c *client) Incr(ctx context.Context, key string) (int64, error) {
	val, err := c.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("incr: %w", err)
	}

	return val, nil
}

func (c *client) Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	ok, err := c.rdb.Expire(ctx, key, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("expire: %w", err)
	}

	return ok, nil
}

func (c *client) Close() error {
	return c.rdb.Close()
}
