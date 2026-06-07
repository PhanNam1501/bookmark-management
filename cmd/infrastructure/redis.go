package infrastructure

import (
	"github.com/redis/go-redis/v9"
	redisPkg "github.com/PhanNam1501/bookmark-management/pkg/redis"
)

// InitRedis initializes Redis client
func InitRedis() (*redis.Client, error) {
	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		return nil, err
	}
	return redisClient, nil
}
