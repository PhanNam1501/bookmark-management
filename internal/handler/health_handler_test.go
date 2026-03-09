package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhanNam1501/bookmark-management/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
)

func TestHealthHandler_GenId(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func() *mocks.Id

		expectedStatus int
		expectResp     string
	}{
		{
			name: "success",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockSvc: func() *mocks.Id {
				svcMock := mocks.NewId(t)
				svcMock.On("GetId").Return("123456789")
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectResp:     "123456789",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)

			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc()
			testHandler := NewHealthCheck(mockSvc)

			testHandler.GenId(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectResp, rec.Body.String())
		})
	}
}
