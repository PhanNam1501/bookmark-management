package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	urlExpTime = 24 * time.Hour
)

//go:generate mockery --name URLStorage --filename urlstorage.go
type URLStorage interface {
	CheckRedisConnection(ctx context.Context) error
	StoreURL(ctx context.Context, code, url string) (string, error)
	LinkShortenURL(ctx context.Context, code, url string, exp int) (string, error)
}

type urlStorage struct {
	redisClient *redis.Client
}

func NewURLStorage(redisClient *redis.Client) URLStorage {
	return &urlStorage{
		redisClient: redisClient,
	}
}

func (s *urlStorage) CheckRedisConnection(ctx context.Context) error {
	pong, err := s.redisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}
	if pong == "PONG" {
		return nil
	}

	return fmt.Errorf("unexpected response from Redis: %s", pong)

}

func (s *urlStorage) StoreURL(ctx context.Context, code, url string) (string, error) {
	//ok, err := s.redisClient.SetNX(ctx, code, url, urlExpTime).Result()
	res, err := s.redisClient.SetArgs(ctx, code, url, redis.SetArgs{
		Mode: "NX", // only set if not exists
		TTL:  urlExpTime,
	}).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s *urlStorage) LinkShortenURL(ctx context.Context, code, url string, exp int) (string, error) {
	res, err := s.redisClient.SetArgs(ctx, code, url, redis.SetArgs{
		Mode: "NX", // only set if not exists
		TTL:  time.Second * time.Duration(exp),
	}).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}
