package infrastructure

import (
	"github.com/PhanNam1501/bookmark-management/internal/api"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// BootstrapDependencies contains all initialized dependencies
type BootstrapDependencies struct {
	Config      *api.Config
	RedisClient *redis.Client
	DB          *gorm.DB
	Engine      api.Engine
}

// Bootstrap initializes all dependencies and returns the app engine
func Bootstrap() (*BootstrapDependencies, error) {
	// 1. Load configuration
	cfg, err := InitConfig()
	if err != nil {
		return nil, err
	}

	// 2. Initialize Redis
	redisClient, err := InitRedis()
	if err != nil {
		return nil, err
	}

	// 3. Initialize Database + Auto Migration
	db, err := InitDatabase()
	if err != nil {
		return nil, err
	}

	// 4. Create Engine/App
	engine := api.New(api.EngineOpts{
		Config:      cfg,
		RedisClient: redisClient,
		DB:          db,
	})

	return &BootstrapDependencies{
		Config:      cfg,
		RedisClient: redisClient,
		DB:          db,
		Engine:      engine,
	}, nil
}
