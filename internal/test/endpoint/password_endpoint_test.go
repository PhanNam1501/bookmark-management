package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhanNam1501/bookmark-management/internal/api"
	redisPkg "github.com/PhanNam1501/bookmark-management/pkg/redis"
	"github.com/go-openapi/testify/v2/assert"
)

func TestPasswordEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHttp func(api.Engine) *httptest.ResponseRecorder

		expectedStatus  int
		expectedRespLen int
	}{
		{
			name: "success",

			setupTestHttp: func(e api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
				respRec := httptest.NewRecorder()
				e.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus:  http.StatusOK,
			expectedRespLen: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg, err := api.NewConfig("")
			if err != nil {
				panic(err)
			}

			redisClient, err := redisPkg.NewClient("")
			if err != nil {
				panic(err)
			}

			app := api.New(cfg, redisClient)
			rec := tc.setupTestHttp(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedRespLen, rec.Body.Len())
		})
	}
}
