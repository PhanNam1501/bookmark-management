package main

import (
	"github.com/PhanNam1501/bookmark-management/cmd/infrastructure"
	"github.com/PhanNam1501/bookmark-management/pkg/logger"
)

// @title			Bookmark API
// @version			1.0
// @description 	Bookmark API
// @host			localhost:8080
// @BasePath		/
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @description JWT token with Bearer prefix
func main() {
	logger.SetLogLevel()

	deps, err := infrastructure.Bootstrap()
	if err != nil {
		panic(err)
	}

	if err := deps.Engine.Start(); err != nil {
		panic(err)
	}
}
