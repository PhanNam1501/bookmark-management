package api

import (
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/handler"
	"github.com/PhanNam1501/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type Engine interface {
	Start() error
	ServeHttp(w http.ResponseWriter, r *http.Request)
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

func (a *api) ServeHttp(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

func (a *api) registerEP() {
	passSvc := service.NewPassword()
	passHandler := handler.NewPassword(passSvc)
	a.app.GET("/gen-pass", passHandler.GenPass)
}
