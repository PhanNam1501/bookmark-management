package endpoint

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhanNam1501/bookmark-management/internal/api"
	redisPkg "github.com/PhanNam1501/bookmark-management/pkg/redis"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/redis/go-redis/v9"
)

func TestURLRedirectEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		seedData func(r *redis.Client)

		setupTestHttp func(api.Engine) *httptest.ResponseRecorder

		expectedStatus   int
		expectedLocation string
	}{
		{
			name: "success",

			seedData: func(r *redis.Client) {
				_ = r.Set(context.Background(), "ScC6OyVRVA", "http://google.com", 0)
			},

			setupTestHttp: func(e api.Engine) *httptest.ResponseRecorder {
				rec := httptest.NewRecorder()

				testCode := "ScC6OyVRVA"

				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/"+testCode, nil)
				e.ServeHTTP(rec, req)
				return rec
			},

			expectedStatus:   http.StatusMovedPermanently,
			expectedLocation: "http://google.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cfg, err := api.NewConfig("")
			if err != nil {
				panic(err)
			}
			redisClient := redisPkg.InitMockRedis(t)
			tc.seedData(redisClient)
			app := api.New(cfg, redisClient)
			rec := tc.setupTestHttp(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedLocation != "" {
				assert.Equal(t, tc.expectedLocation, rec.Header().Get("Location"))
			}
		})
	}
}
