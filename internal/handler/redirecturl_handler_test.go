package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhanNam1501/bookmark-management/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/stretchr/testify/mock"
)

func TestRedirectURLHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func() (*gin.Context, *httptest.ResponseRecorder)
		setupMockSvc func() *mocks.URLRedirect

		expectedStatus int
	}{
		{
			name: "success",
			setupRequest: func() (*gin.Context, *httptest.ResponseRecorder) {
				rec := httptest.NewRecorder()
				gc, _ := gin.CreateTestContext(rec)

				testCode := "ScC6OyVRVA"
				gc.Params = []gin.Param{{Key: "code", Value: testCode}}
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/"+testCode, nil)
				req.Header.Set("Content-Type", "application/json")

				gc.Request = req
				return gc, rec
			},
			setupMockSvc: func() *mocks.URLRedirect {
				svcMock := mocks.NewURLRedirect(t)
				svcMock.On(
					"GetRedirectURL",
					mock.Anything,
					"ScC6OyVRVA",
				).Return("http://google.com", nil)
				return svcMock
			},

			expectedStatus: http.StatusMovedPermanently,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gc, rec := tc.setupRequest()
			mockSvc := tc.setupMockSvc()
			testHandler := NewRedirectURLHandler(mockSvc)

			testHandler.RedirectURL(gc)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
