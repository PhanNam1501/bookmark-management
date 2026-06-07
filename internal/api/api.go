package api

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/PhanNam1501/bookmark-management/docs"
	"github.com/PhanNam1501/bookmark-management/internal/api/middlewares"
	"github.com/PhanNam1501/bookmark-management/internal/handler"
	bookmarkhandler "github.com/PhanNam1501/bookmark-management/internal/handler/bookmark"
	"github.com/PhanNam1501/bookmark-management/internal/repository"
	"github.com/PhanNam1501/bookmark-management/internal/repository/cache"
	"github.com/PhanNam1501/bookmark-management/internal/repository/ratelimit"
	"github.com/PhanNam1501/bookmark-management/internal/service"
	"github.com/PhanNam1501/bookmark-management/internal/service/bookmark"
	"github.com/PhanNam1501/bookmark-management/pkg/jwtutils"
	"github.com/PhanNam1501/bookmark-management/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type EngineOpts struct {
	Config      *Config
	RedisClient *redis.Client
	DB          *gorm.DB
}

type api struct {
	app          *gin.Engine
	cfg          *Config
	redisClient  *redis.Client
	db           *gorm.DB
	jwtValidator jwtutils.JWTValidator
}

func New(opts EngineOpts) Engine {
	jwtValidator, err := jwtutils.NewJWTValidator(filepath.FromSlash("./public.pem"))
	if err != nil {
		panic("Failed to initialize JWT validator: " + err.Error())
	}

	a := &api{
		app:          gin.Default(),
		cfg:          opts.Config,
		redisClient:  opts.RedisClient,
		db:           opts.DB,
		jwtValidator: jwtValidator,
	}

	a.registerRoutes()

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

type AppMiddlewares struct {
	jwtAuth   middlewares.JWTAuth
	rateLimit middlewares.RateLimit
}

func (a *api) initMiddlewares() AppMiddlewares {
	rateLimitRepo := ratelimit.NewRedisRepository(a.redisClient)
	return AppMiddlewares{
		jwtAuth:   middlewares.NewJWTAuth(a.jwtValidator),
		rateLimit: middlewares.NewRateLimit(rateLimitRepo),
	}
}

func (a *api) registerRoutes() {
	handlers := a.getHandlers()

	// Public routes
	a.app.GET("/gen-pass", handlers.Password.GenPass)
	a.app.GET("/health-check", handlers.Bookmark.GenUuid)
	a.app.POST("/shorten", handlers.ShortenURL.ShortenURL)
	a.app.POST("/v1/links/shorten", handlers.LinkShorten.LinkShortenURL)
	a.app.GET("/v1/links/redirect/:code", handlers.RedirectURL.RedirectURL)

	a.app.POST("/users/register", handlers.User.RegisterUser)
	a.app.POST("/users/login", handlers.User.Login)
	a.app.POST("/test/:code", handlers.User.Test)

	// Protected routes (require JWT)
	privateRoutes := a.app.Group("")
	mw := a.initMiddlewares()
	privateRoutes.Use(mw.jwtAuth.JWTAuth())
	privateRoutes.Use(mw.rateLimit.RateLimit())
	privateRoutes.GET("/v1/self/info", handlers.User.GetCurrentUser)
	privateRoutes.POST("/v1/bookmarks", handlers.BookmarkHandler.CreateBookmark)
	privateRoutes.GET("/v1/bookmarks", handlers.BookmarkHandler.GetBookmarks)
	privateRoutes.PUT("/v1/bookmarks/:id", handlers.BookmarkHandler.UpdateBookmark)
	privateRoutes.DELETE("/v1/bookmarks/:id", handlers.BookmarkHandler.DeleteBookmark)
}

type Handlers struct {
	Password        handler.Password
	Bookmark        handler.Bookmark
	BookmarkHandler bookmarkhandler.Handler
	ShortenURL      handler.ShortenURL
	LinkShorten     handler.LinkShortURL
	RedirectURL     handler.RedirectURL
	User            handler.User
}

func (a *api) getHandlers() *Handlers {
	// Password handler
	passSvc := service.NewPassword()
	passHandler := handler.NewPasswordHandler(passSvc)

	// Bookmark handler
	urlStorage := repository.NewURLStorage(a.redisClient)
	bookmarkSvc := service.NewBookmark(urlStorage)
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkSvc)

	// Shorten URL handler
	passwordSvc := service.NewPassword()
	shortenURLSvc := service.NewShortenURL(urlStorage, passwordSvc)
	shortenURLHandler := handler.NewShortenURLHandler(shortenURLSvc)

	// Link Shorten handler
	linkShortenHandler := handler.NewLinkShortenHandler(shortenURLSvc)

	// Bookmark handler (with DB) - needed for redirect service
	bookmarkRepo := repository.NewRepository(a.db)

	// Redirect URL handler - support both Redis (code len 7) and DB (code len 8)
	urlRedirectSvc := service.NewUrlRedirectWithDB(urlStorage, bookmarkRepo)
	redirectURLHandler := handler.NewRedirectURLHandler(urlRedirectSvc)

	cacheRepo := cache.NewRedisRepo(a.redisClient)
	keyGen := utils.NewRandomKeyGenerator(16)
	bookmarkSvcDB := bookmark.NewService(bookmarkRepo, keyGen)
	bookmarkSvcWithCache := bookmark.NewServiceWithCache(bookmarkSvcDB, cacheRepo)
	bookmarkHandlerDB := bookmarkhandler.NewHandler(bookmarkSvcWithCache)

	// User handler
	userRepo := repository.NewUser(a.db)
	jwtGen, err := jwtutils.NewJWTGenerator(filepath.FromSlash("./private.pem"))
	if err != nil {
		panic("Failed to initialize JWT generator: " + err.Error())
	}
	userSvc := service.NewUser(userRepo, jwtGen)
	userHandler := handler.NewUser(userSvc)

	return &Handlers{
		Password:        passHandler,
		Bookmark:        bookmarkHandler,
		BookmarkHandler: bookmarkHandlerDB,
		ShortenURL:      shortenURLHandler,
		LinkShorten:     linkShortenHandler,
		RedirectURL:     redirectURLHandler,
		User:            userHandler,
	}
}
