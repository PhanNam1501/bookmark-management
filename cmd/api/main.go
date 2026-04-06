package main

import (
	"github.com/PhanNam1501/bookmark-management/internal/api"
	"github.com/PhanNam1501/bookmark-management/pkg/logger"
	redisPkg "github.com/PhanNam1501/bookmark-management/pkg/redis"
)

// @title			Bookmark API
// @version			1.0
// @description 	Bookmark API
// @host			localhost:8080
// @BasePath		/
func main() {
	logger.SetLogLevel()

	cfg, err := api.NewConfig("")
	if err != nil {
		panic(err)
	}
	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}
	app := api.New(cfg, redisClient)
	app.Start()
}
