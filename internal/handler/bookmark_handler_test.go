package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhanNam1501/bookmark-management/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
)

func TestBookmarkHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func() (*gin.Context, *httptest.ResponseRecorder)
		setupMockSvc func() *mocks.Bookmark

		expectedStatus int
		expectedResp   BaseResponse
	}{
		{
			name: "success",
			setupRequest: func() (*gin.Context, *httptest.ResponseRecorder) {
				rec := httptest.NewRecorder()
				gc, _ := gin.CreateTestContext(rec)
				gc.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
				return gc, rec
			},

			setupMockSvc: func() *mocks.Bookmark {
				svcMock := mocks.NewBookmark(t)
				svcMock.On("GenerateUuid").Return("15012004") //have input -> , , , ,
				return svcMock
			},

			expectedStatus: http.StatusOK,
			expectedResp: BaseResponse{
				Message:     "OK",
				ServiceName: "bookmark_service",
				InstanceID:  "15012004",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gc, rec := tc.setupRequest()
			mockSvc := tc.setupMockSvc()
			testHandler := NewBookmarkHandler(mockSvc)

			testHandler.GenUuid(gc)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			var resp BaseResponse
			json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.Equal(t, tc.expectedResp, resp)
		})
	}
}
