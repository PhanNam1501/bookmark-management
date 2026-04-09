package api

import (
	"fmt"
	"net/http"

	"github.com/PhanNam1501/bookmark-management/docs"
	"github.com/PhanNam1501/bookmark-management/internal/handler"
	"github.com/PhanNam1501/bookmark-management/internal/repository"
	"github.com/PhanNam1501/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Engine interface {
	Start() error
	registerEP()
	createUUID()
	shortenURL()
	linkShortenURL()
	redirectURL()
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type api struct {
	app         *gin.Engine
	cfg         *Config
	redisClient *redis.Client
}

func New(cfg *Config, redisClient *redis.Client) Engine {
	a := &api{
		app:         gin.Default(),
		cfg:         cfg,
		redisClient: redisClient,
	}
	a.registerEP()
	a.createUUID()
	a.shortenURL()
	a.linkShortenURL()
	a.redirectURL()

	docs.SwaggerInfo.Host = a.cfg.Hostname
	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
	urlStorage := repository.NewURLStorage(a.redisClient)
	bookmarkSvc := service.NewBookmark(urlStorage)
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkSvc)
	a.app.GET("/health-check", bookmarkHandler.GenUuid)
}

func (a *api) shortenURL() {
	passwordSvc := service.NewPassword()
	urlStorage := repository.NewURLStorage(a.redisClient)
	shortenURLSvc := service.NewShortenURL(urlStorage, passwordSvc)
	shortenURLHandler := handler.NewShortenURLHandler(shortenURLSvc)

	a.app.POST("/shorten", shortenURLHandler.ShortenURL)
}

func (a *api) linkShortenURL() {
	passwordSvc := service.NewPassword()
	urlStorage := repository.NewURLStorage(a.redisClient)
	shortenURLSvc := service.NewShortenURL(urlStorage, passwordSvc)
	shortenURLHandler := handler.NewLinkShortenHandler(shortenURLSvc)

	a.app.POST("/v1/links/shorten", shortenURLHandler.LinkShortenURL)
}

func (a *api) redirectURL() {
	urlStorage := repository.NewURLStorage(a.redisClient)
	urlRedirectSvc := service.NewUrlRedirect(urlStorage)
	urlRedirectHandler := handler.NewRedirectURLHandler(urlRedirectSvc)

	a.app.GET("/v1/links/redirect/:code", urlRedirectHandler.RedirectURL)
}
