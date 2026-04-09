package endpoint

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhanNam1501/bookmark-management/internal/api"
	"github.com/PhanNam1501/bookmark-management/internal/handler"
	redisPkg "github.com/PhanNam1501/bookmark-management/pkg/redis"
	"github.com/go-openapi/testify/v2/assert"
)

func TestBookmarkEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHttp func(api.Engine) *httptest.ResponseRecorder

		expectedStatus          int
		expectedRespMessage     string
		expectedRespServiceName string
		expectedRespLen         int
	}{
		{
			name: "success",

			setupTestHttp: func(e api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
				respRec := httptest.NewRecorder()
				e.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus:          http.StatusOK,
			expectedRespMessage:     "OK",
			expectedRespServiceName: "bookmark_service",
			expectedRespLen:         36,
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
			var resp handler.BaseResponse
			json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.Equal(t, tc.expectedRespMessage, resp.Message)
			assert.Equal(t, tc.expectedRespServiceName, resp.ServiceName)
			assert.Equal(t, tc.expectedRespLen, len(resp.InstanceID))
		})
	}
}
