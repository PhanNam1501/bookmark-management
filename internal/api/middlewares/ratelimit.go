package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PhanNam1501/bookmark-management/internal/handler/utils"
	"github.com/PhanNam1501/bookmark-management/internal/repository/ratelimit"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	rateLimitFormat  = "rate_limit:%s"
	rateLimitExpTime = 1 * time.Minute
	rateLimitMaxRate = 10
)

type RateLimit interface {
	RateLimit() gin.HandlerFunc
}

type rateLimitMiddleware struct {
	repo ratelimit.Repository
}

func NewRateLimit(repo ratelimit.Repository) RateLimit {
	return &rateLimitMiddleware{
		repo: repo,
	}
}

func (r *rateLimitMiddleware) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get key from uid
		uid, err := utils.GetUserIDFromRequest(c)
		if err != nil {
			return
		}
		rateLimitKey := fmt.Sprintf(rateLimitFormat, uid)
		// check rate limit
		curRateLimit, err := r.repo.GetRateLimit(c, rateLimitKey)
		if err != nil {
			log.Error().Err(err).Str("uid", uid).Str("key", rateLimitKey).Msg("get rate limit")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to check rate limit"})
			return
		}

		if curRateLimit >= rateLimitMaxRate {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		if err := r.repo.IncreaseRateLimit(c, rateLimitKey, rateLimitExpTime); err != nil {
			log.Error().Err(err).Str("uid", uid).Str("key", rateLimitKey).Msg("increase rate limit")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rate limit"})
			return
		}

		// call next
		c.Next()
	}
}
