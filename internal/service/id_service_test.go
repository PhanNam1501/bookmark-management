package service

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestIdService_GenerateId(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		expectedLen int
	}{
		{
			name: "normal case",

			expectedLen: 36,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testSvc := NewId()
			id := testSvc.GetId()

			assert.Equal(t, tc.expectedLen, len(id))
		})
	}
}
