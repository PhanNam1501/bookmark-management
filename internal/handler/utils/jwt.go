package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrEmptyUID     = errors.New("empty uid")
)

const (
	jwtClaimsKey = "claims"
)

// SetJWTClaims stores JWT claims in the context
func SetJWTClaims(c *gin.Context, claims jwt.MapClaims) {
	c.Set(jwtClaimsKey, claims)
}

// GetJWTClaimsFromRequest returns jwt.MapClaims from request context
func GetJWTClaimsFromRequest(c *gin.Context) (jwt.MapClaims, error) {
	tokenInfo, exists := c.Get(jwtClaimsKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, BindError{
			Code:   "INVALID_TOKEN",
			Errors: []InputFieldError{{Field: "token", Message: ErrInvalidToken.Error()}},
		})
		return nil, ErrInvalidToken
	}

	claims, valid := tokenInfo.(jwt.MapClaims)
	if !valid {
		c.JSON(http.StatusUnauthorized, BindError{
			Code:   "INVALID_TOKEN",
			Errors: []InputFieldError{{Field: "token", Message: ErrInvalidToken.Error()}},
		})
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetUserIDFromRequest returns user id from request's token
func GetUserIDFromRequest(c *gin.Context) (string, error) {
	claims, err := GetJWTClaimsFromRequest(c)
	if err != nil {
		return "", err
	}

	uid, exists := claims["sub"]
	if !exists {
		c.JSON(http.StatusUnauthorized, BindError{
			Code:   "EMPTY_UID",
			Errors: []InputFieldError{{Field: "uid", Message: ErrEmptyUID.Error()}},
		})
		return "", ErrEmptyUID
	}

	uidStr, ok := uid.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, BindError{
			Code:   "INVALID_UID",
			Errors: []InputFieldError{{Field: "uid", Message: "uid is not a string"}},
		})
		return "", errors.New("uid is not a string")
	}

	if uidStr == "" {
		c.JSON(http.StatusUnauthorized, BindError{
			Code:   "EMPTY_UID",
			Errors: []InputFieldError{{Field: "uid", Message: ErrEmptyUID.Error()}},
		})
		return "", ErrEmptyUID
	}

	return uidStr, nil
}

// GetClaimValue retrieves a claim value by key
func GetClaimValue[T any](c *gin.Context, key string) (T, error) {
	var zero T
	claims, err := GetJWTClaimsFromRequest(c)
	if err != nil {
		return zero, err
	}

	value, exists := claims[key]
	if !exists {
		return zero, errors.New("claim key not found")
	}

	typedValue, ok := value.(T)
	if !ok {
		return zero, errors.New("claim value type mismatch")
	}

	return typedValue, nil
}
