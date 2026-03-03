package main

import "github.com/PhanNam1501/bookmark-management/internal/api"

func main() {
	app := api.New()
	app.Start()

	_ = api.NewMux().StartMux()
}
