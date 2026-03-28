package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PhanNam1501/bookmark-management/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenURLHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func() (*gin.Context, *httptest.ResponseRecorder)
		setupMockSvc func() *mocks.ShortenURL

		expectedStatus int
		expectedLen    int
		expectedResp   ShortenURLResponse
	}{
		{
			name: "success",
			setupRequest: func() (*gin.Context, *httptest.ResponseRecorder) {
				rec := httptest.NewRecorder()
				gc, _ := gin.CreateTestContext(rec)

				body := `{
					"url": "http://google.com"
				}`

				req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				gc.Request = req
				return gc, rec
			},
			setupMockSvc: func() *mocks.ShortenURL {
				svcMock := mocks.NewShortenURL(t)
				svcMock.On(
					"ShortenURL",
					mock.Anything,
					"http://google.com",
				).Return("1234567890", nil)
				return svcMock
			},

			expectedStatus: http.StatusOK,
			expectedLen:    10,
			expectedResp: ShortenURLResponse{
				Key: "1234567890",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gc, rec := tc.setupRequest()
			mockSvc := tc.setupMockSvc()
			testHandler := NewShortenURLHandler(mockSvc)

			testHandler.ShortenURL(gc)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			var resp ShortenURLResponse
			json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.Equal(t, tc.expectedResp, resp)
		})
	}
}
