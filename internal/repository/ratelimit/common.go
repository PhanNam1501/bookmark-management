package ratelimit

import (
	"context"
	"time"
)

//go:generate mockery --name Repository --filename repository_mock.go

type Repository interface {
	IncreaseRateLimit(ctx context.Context, key string, exp time.Duration) error
	GetRateLimit(ctx context.Context, key string) (int, error)
}
