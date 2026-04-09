package service

import (
	"testing"

	"github.com/PhanNam1501/bookmark-management/internal/repository/mocks"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/stretchr/testify/mock"
)

func TestURLRedirectService_GetRedirectURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockSvc func() *mocks.URLStorage

		expectedURL string
		expectedErr error
	}{
		{
			name: "normal case",

			setupMockSvc: func() *mocks.URLStorage {
				svcMock := mocks.NewURLStorage(t)
				svcMock.On(
					"GetRedirectURL",
					mock.Anything,
					"ScC6OyVRVA",
				).Return("https://www.google.com/", nil)

				return svcMock
			},

			expectedURL: "https://www.google.com/",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockUrl := tc.setupMockSvc()
			testSvc := NewUrlRedirect(mockUrl)
			url, err := testSvc.GetRedirectURL(t.Context(), "ScC6OyVRVA")

			assert.Equal(t, tc.expectedURL, url)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
