package ratelimit

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	c *redis.Client
}

func NewRedisRepository(c *redis.Client) Repository {
	return &redisRepository{c: c}
}

func (r *redisRepository) IncreaseRateLimit(ctx context.Context, key string, exp time.Duration) error {
	_, err := r.c.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		err := r.c.Incr(ctx, key).Err()
		if err != nil {
			return err
		}
		err = r.c.Expire(ctx, key, exp).Err()
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *redisRepository) GetRateLimit(ctx context.Context, key string) (int, error) {
	curRateLimit, err := r.c.Get(ctx, key).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}

		return -1, err
	}

	return curRateLimit, nil
}
