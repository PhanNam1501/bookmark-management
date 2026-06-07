package cache

//go:generate mockery --name Repo --filename repo_mock.go

import (
	"context"
	"time"
)

type Repo interface {
	SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value []byte, exp time.Duration) error
	GetCacheData(ctx context.Context, cacheGroupKey, cacheKey string) ([]byte, error)
	DeleteCacheData(ctx context.Context, cacheGroupKey string) error
}
