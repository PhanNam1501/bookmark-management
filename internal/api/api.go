package api

import (
	"github.com/PhanNam1501/bookmark-management/internal/handler"
	"github.com/PhanNam1501/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type Engine interface {
	Start() error
}

type api struct {
	app *gin.Engine
}

func New() Engine {
	a := &api{
		app: gin.New(),
	}
	a.registerEP()
	return a
}

func (a *api) Start() error {
	return a.app.Run(":8080")
}

func (a *api) registerEP() {
	passSvc := service.NewPassword()
	passHandler := handler.NewPassword(passSvc)
	a.app.GET("/gen-pass", passHandler.GenPass)
}
