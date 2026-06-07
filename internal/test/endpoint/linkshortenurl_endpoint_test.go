package endpoint

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PhanNam1501/bookmark-management/internal/api"
	"github.com/PhanNam1501/bookmark-management/internal/test/fixture"
	redisPkg "github.com/PhanNam1501/bookmark-management/pkg/redis"
	"github.com/go-openapi/testify/v2/assert"
)

type LinkShortenURLResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func TestLinkShortenURLEndpoint(t *testing.T) {
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
				body := `{
					"exp": 64000,
					"url": "http://google.com"
				}`
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(body))
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

			redisClient := redisPkg.InitMockRedis(t)
			db := fixture.NewFixture(t, &fixture.UserCommonTestDB{})

			app := api.New(api.EngineOpts{
				Config:      cfg,
				RedisClient: redisClient,
				DB:          db,
			})
			rec := tc.setupTestHttp(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			var resp LinkShortenURLResponse
			json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.Equal(t, tc.expectedRespLen, len(resp.Code))
		})
	}
}
