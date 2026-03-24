package service

import (
	"testing"

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

			testSvc := NewBookmark()
			uuid := testSvc.GenerateUuid()

			assert.Equal(t, tc.expectedLen, len(uuid))
		})
	}
}
