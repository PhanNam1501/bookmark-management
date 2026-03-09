package api

import (
	"fmt"
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/handler"
	"github.com/PhanNam1501/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type api struct {
	app *gin.Engine
	cfg *Config
}

func New(cfg *Config) Engine {
	a := &api{
		app: gin.New(),
		cfg: cfg,
	}
	a.registerEP()
	a.getId()
	return a
}

func (a *api) Start() error {
	return a.app.Run(fmt.Sprintf(":%s", a.cfg.AppPort))
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

func (a *api) registerEP() {
	passSvc := service.NewPassword()
	passHandler := handler.NewPassword(passSvc)
	a.app.GET("/gen-pass", passHandler.GenPass)
}

func (a *api) getId() {
	idSvc := service.NewId()
	idHander := handler.NewHealthCheck(idSvc)
	a.app.GET("/health-check", idHander.GenId)
}
