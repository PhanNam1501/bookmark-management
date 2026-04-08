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

func TestLinkShortenURLHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func() (*gin.Context, *httptest.ResponseRecorder)
		setupMockSvc func() *mocks.ShortenURL

		expectedStatus int
		expectedLen    int
		expectedResp   LinkShortenURLResponse
	}{
		{
			name: "success",
			setupRequest: func() (*gin.Context, *httptest.ResponseRecorder) {
				rec := httptest.NewRecorder()
				gc, _ := gin.CreateTestContext(rec)

				body := `{
					"exp": 64000,
					"url": "http://google.com"
				}`

				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				gc.Request = req
				return gc, rec
			},
			setupMockSvc: func() *mocks.ShortenURL {
				svcMock := mocks.NewShortenURL(t)
				svcMock.On(
					"LinkShortenURL",
					mock.Anything,
					"http://google.com",
					64000,
				).Return("1234567890", nil)
				return svcMock
			},

			expectedStatus: http.StatusOK,
			expectedLen:    10,
			expectedResp: LinkShortenURLResponse{
				Code:    "1234567890",
				Message: "Shorten URL generated successfully!",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gc, rec := tc.setupRequest()
			mockSvc := tc.setupMockSvc()
			testHandler := NewLinkShortenHandler(mockSvc)

			testHandler.LinkShortenURL(gc)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			var resp LinkShortenURLResponse
			json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.Equal(t, tc.expectedResp, resp)
		})
	}
}
