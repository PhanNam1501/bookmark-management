package main

import (
	"context"
	"time"

	"github.com/PhanNam1501/bookmark-management/pkg/redis"
)

func main() {
	redisClient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}

	redisClient.Set(context.Background(), "test", "test", time.Hour)
}
