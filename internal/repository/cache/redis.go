package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisRepo struct {
	redisClient *redis.Client
}

func NewRedisRepo(r *redis.Client) Repo {
	return &redisRepo{
		redisClient: r,
	}
}

func (r *redisRepo) SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value []byte, exp time.Duration) error {
	err := r.redisClient.HSet(ctx, cacheGroupKey, cacheKey, value).Err()
	if err != nil {
		return err
	}

	return r.redisClient.Expire(ctx, cacheGroupKey, exp).Err()
}

func (r *redisRepo) GetCacheData(ctx context.Context, cacheGroupKey, cacheKey string) ([]byte, error) {
	val, err := r.redisClient.HGet(ctx, cacheGroupKey, cacheKey).Bytes()
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (r *redisRepo) DeleteCacheData(ctx context.Context, cacheGroupKey string) error {
	err := r.redisClient.Del(ctx, cacheGroupKey).Err()
	if err != nil {
		return err
	}
	return nil
}
