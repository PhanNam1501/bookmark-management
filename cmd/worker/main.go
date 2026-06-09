package main

import (
	"context"

	"github.com/PhanNam1501/bookmark-management/internal/handler/workerhandler"
	"github.com/PhanNam1501/bookmark-management/internal/repository/queue"
	"github.com/PhanNam1501/bookmark-management/internal/worker"
	"github.com/PhanNam1501/bookmark-management/pkg/redis"
)

func main() {

	ctx := context.Background()
	// init redis
	redisClient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}
	// init queue repo
	queueRepo := queue.NewRedisQueue(redisClient, "bookmark_queue")

	// init handler
	testHandler := &workerhandler.TestHandler{}
	// init engine
	workerEngine := worker.NewEngine(queueRepo, testHandler)
	// start engine
	workerEngine.Start(ctx)
}
