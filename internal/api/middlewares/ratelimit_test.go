package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PhanNam1501/bookmark-management/internal/handler/utils"
	"github.com/PhanNam1501/bookmark-management/internal/repository/ratelimit/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestContext(uid string) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	claims := jwt.MapClaims{
		"sub": uid,
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	utils.SetJWTClaims(c, claims)
	return c
}

func TestRateLimit_RateLimit_Success(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.Repository)
	testKey := fmt.Sprintf(rateLimitFormat, "user123")

	mockRepo.On("GetRateLimit", mock.Anything, testKey).Return(5, nil)
	mockRepo.On("IncreaseRateLimit", mock.Anything, testKey, rateLimitExpTime).Return(nil)

	middleware := NewRateLimit(mockRepo)
	handler := middleware.RateLimit()

	c := setupTestContext("user123")
	handler(c)

	assert.False(t, c.IsAborted(), "request should not be aborted when under limit")
	mockRepo.AssertExpectations(t)
}

func TestRateLimit_RateLimit_ExceededLimit(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.Repository)
	testKey := fmt.Sprintf(rateLimitFormat, "user456")

	mockRepo.On("GetRateLimit", mock.Anything, testKey).Return(rateLimitMaxRate, nil)

	middleware := NewRateLimit(mockRepo)
	handler := middleware.RateLimit()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	claims := jwt.MapClaims{
		"sub": "user456",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	utils.SetJWTClaims(c, claims)

	handler(c)

	assert.True(t, c.IsAborted(), "request should be aborted when limit exceeded")
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	mockRepo.AssertNotCalled(t, "IncreaseRateLimit")
}

func TestRateLimit_RateLimit_GetRateLimitError(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.Repository)
	testKey := fmt.Sprintf(rateLimitFormat, "user789")
	testError := errors.New("redis connection error")

	mockRepo.On("GetRateLimit", mock.Anything, testKey).Return(-1, testError)

	middleware := NewRateLimit(mockRepo)
	handler := middleware.RateLimit()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	claims := jwt.MapClaims{
		"sub": "user789",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	utils.SetJWTClaims(c, claims)

	handler(c)

	assert.True(t, c.IsAborted(), "request should be aborted on GetRateLimit error")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertNotCalled(t, "IncreaseRateLimit")
}

func TestRateLimit_RateLimit_IncreaseRateLimitError(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.Repository)
	testKey := fmt.Sprintf(rateLimitFormat, "user999")
	testError := errors.New("redis write error")

	mockRepo.On("GetRateLimit", mock.Anything, testKey).Return(3, nil)
	mockRepo.On("IncreaseRateLimit", mock.Anything, testKey, rateLimitExpTime).Return(testError)

	middleware := NewRateLimit(mockRepo)
	handler := middleware.RateLimit()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	claims := jwt.MapClaims{
		"sub": "user999",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	utils.SetJWTClaims(c, claims)

	handler(c)

	assert.True(t, c.IsAborted(), "request should be aborted on IncreaseRateLimit error")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestRateLimit_RateLimit_NoUserID(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.Repository)

	middleware := NewRateLimit(mockRepo)
	handler := middleware.RateLimit()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	// không set JWT claims

	handler(c)

	// Nếu không có uid, middleware sẽ abort (không call repo)
	mockRepo.AssertNotCalled(t, "GetRateLimit")
	mockRepo.AssertNotCalled(t, "IncreaseRateLimit")
}

func TestRateLimit_RateLimit_AtLimit(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.Repository)
	testKey := fmt.Sprintf(rateLimitFormat, "user_at_limit")

	// Test with rateLimitMaxRate - 1 (just under limit)
	mockRepo.On("GetRateLimit", mock.Anything, testKey).Return(rateLimitMaxRate - 1, nil)
	mockRepo.On("IncreaseRateLimit", mock.Anything, testKey, rateLimitExpTime).Return(nil)

	middleware := NewRateLimit(mockRepo)
	handler := middleware.RateLimit()

	c := setupTestContext("user_at_limit")
	handler(c)

	assert.False(t, c.IsAborted(), "request should pass when at limit-1")
	mockRepo.AssertExpectations(t)
}
