package repository

import (
	"context"
	"testing"

	redisPkg "github.com/PhanNam1501/bookmark-management/pkg/redis"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/redis/go-redis/v9"
)

func TestURLStoratge(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func() *redis.Client
		inputCode string
		inputURL  string

		expectedErr error
		verifyFunc  func(ctx context.Context, r *redis.Client, inputCode, inputURl string, ok string)
	}{
		{
			name: "normal case",

			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			inputCode:   "12345",
			inputURL:    "http://google.com",
			expectedErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client, inputCode, inputURl string, ok string) {
				res, err := r.Get(ctx, inputCode).Result()
				assert.NoError(t, err)
				assert.Equal(t, inputURl, res)
				assert.Equal(t, ok, "OK")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			redisClient := tc.setupMock()
			testRepo := NewURLStorage(redisClient)

			ok, err := testRepo.StoreURL(ctx, tc.inputCode, tc.inputURL)
			if err == nil {
				tc.verifyFunc(ctx, redisClient, tc.inputCode, tc.inputURL, ok)
			}
		})
	}
}
