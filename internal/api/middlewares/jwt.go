package middlewares

import (
	"net/http"
	"strings"

	"github.com/PhanNam1501/bookmark-management/pkg/jwtutils"
	"github.com/gin-gonic/gin"
)

type JWTAuth interface {
	JWTAuth() gin.HandlerFunc
}

type jwtAuth struct {
	jwtValidator jwtutils.JWTValidator
}

func NewJWTAuth(jwtValidator jwtutils.JWTValidator) JWTAuth {
	return &jwtAuth{
		jwtValidator: jwtValidator,
	}
}

func (j *jwtAuth) JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get auth header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Auth header is required"})
			ctx.Abort()
			return
		}
		// get token from auth header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth header formal"})
			ctx.Abort()
			return
		}
		tokenString := parts[1]
		// validate token
		claims, err := j.jwtValidator.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// set claims to context
		ctx.Set("claims", claims)

		ctx.Next()
	}
}
