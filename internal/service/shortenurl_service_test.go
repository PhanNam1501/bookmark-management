package service

import (
	"testing"

	"github.com/PhanNam1501/bookmark-management/internal/repository/mocks"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenURLService_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		url string

		setupMockSvc func() *mocks.URLStorage

		expectedLen int
		expectedErr error
	}{
		{
			name: "normal case",
			url:  "http://google.com",
			setupMockSvc: func() *mocks.URLStorage {
				svcMock := mocks.NewURLStorage(t)
				svcMock.On(
					"StoreURL",
					mock.Anything,
					mock.AnythingOfType("string"),
					"http://google.com",
				).Return("OK", nil)

				return svcMock
			},
			expectedLen: 10,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockUrl := tc.setupMockSvc()
			svcPassword := NewPassword()
			testSvc := NewShortenURL(mockUrl, svcPassword)
			code, err := testSvc.ShortenURL(t.Context(), tc.url)

			assert.Equal(t, tc.expectedLen, len(code))
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, urlSafeRegex.MatchString(code), true)
		})
	}
}

func TestLinkShortenURLService_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		url string
		exp int

		setupMockSvc func() *mocks.URLStorage

		expectedLen int
		expectedErr error
	}{
		{
			name: "normal case",
			url:  "http://google.com",
			exp:  64000,
			setupMockSvc: func() *mocks.URLStorage {
				svcMock := mocks.NewURLStorage(t)
				svcMock.On(
					"LinkShortenURL",
					mock.Anything,
					mock.AnythingOfType("string"),
					"http://google.com",
					64000,
				).Return("OK", nil)

				return svcMock
			},
			expectedLen: 10,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockUrl := tc.setupMockSvc()
			svcPassword := NewPassword()
			testSvc := NewShortenURL(mockUrl, svcPassword)
			code, err := testSvc.LinkShortenURL(t.Context(), tc.url, tc.exp)

			assert.Equal(t, tc.expectedLen, len(code))
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, urlSafeRegex.MatchString(code), true)
		})
	}
}
