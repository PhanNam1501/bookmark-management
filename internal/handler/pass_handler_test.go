package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhanNam1501/bookmark-management/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
)

func TestPasswordHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func() (*gin.Context, *httptest.ResponseRecorder)
		setupMockSvc func() *mocks.Password

		expectedStatus int
		expectedResp   string
	}{
		{
			name: "success",
			// clean
			setupRequest: func() (*gin.Context, *httptest.ResponseRecorder) {
				rec := httptest.NewRecorder()
				gc, _ := gin.CreateTestContext(rec)
				gc.Request = httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
				return gc, rec
			},
			setupMockSvc: func() *mocks.Password {
				svcMock := mocks.NewPassword(t)
				svcMock.On("GeneratePassword").Return("123456789", nil) //have input -> , , , ,
				return svcMock
			},

			expectedStatus: http.StatusOK,
			expectedResp:   "123456789",
		},

		{
			name: "internal server err",

			setupRequest: func() (*gin.Context, *httptest.ResponseRecorder) {
				rec := httptest.NewRecorder()
				gc, _ := gin.CreateTestContext(rec)
				gc.Request = httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
				return gc, rec
			},
			setupMockSvc: func() *mocks.Password {
				svcMock := mocks.NewPassword(t)
				svcMock.On("GeneratePassword").Return("", errors.New("something")) //have input -> , , , ,
				return svcMock
			},

			expectedStatus: http.StatusInternalServerError,
			expectedResp:   "err",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// fake responsewriter
			// rec := httptest.NewRecorder()
			// gc, _ := gin.CreateTestContext(rec)
			gc, rec := tc.setupRequest()
			mockSvc := tc.setupMockSvc()
			testHandler := NewPasswordHandler(mockSvc)

			testHandler.GenPass(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResp, rec.Body.String())
		})
	}
}
