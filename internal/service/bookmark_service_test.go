package service

import (
	"testing"

	"github.com/PhanNam1501/bookmark-management/internal/repository/mocks"
	"github.com/go-openapi/testify/v2/assert"
)

func TestBookmarkService(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		expectedLen int
	}{
		{
			name:        "normal case",
			expectedLen: 36,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockUrl := mocks.NewURLStorage(t)
			testSvc := NewBookmark(mockUrl)
			uuid := testSvc.GenerateUuid()

			assert.Equal(t, tc.expectedLen, len(uuid))
		})
	}
}
