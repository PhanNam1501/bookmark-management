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
	registerEP()
	createUUID()
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
	a.createUUID()
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
	passHandler := handler.NewPasswordHandler(passSvc)
	a.app.GET("/gen-pass", passHandler.GenPass)
}

func (a *api) createUUID() {
	bookmarkSvc := service.NewBookmark()
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkSvc)
	a.app.GET("/health-check", bookmarkHandler.GenUuid)
}
