package main

import "github.com/PhanNam1501/bookmark-management/internal/api"

// @title  Bookmark API
// version 1.0
// @description Bookmark API
// @host 		localhost:8080
// @BasePath    /
func main() {
	cfg, err := api.NewConfig("")
	if err != nil {
		panic(err)
	}

	app := api.New(cfg)
	app.Start()
}
