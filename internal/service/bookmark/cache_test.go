package bookmark_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/PhanNam1501/bookmark-management/internal/model"
	cacheMocks "github.com/PhanNam1501/bookmark-management/internal/repository/cache/mocks"
	bookmark "github.com/PhanNam1501/bookmark-management/internal/service/bookmark"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBookmarkService struct {
	mock.Mock
}

func (m *mockBookmarkService) CreateBookmark(ctx context.Context, userid, url, description string) (*model.Bookmark, error) {
	args := m.Called(ctx, userid, url, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bookmark), args.Error(1)
}

func (m *mockBookmarkService) GetBookmarks(ctx context.Context, userID string, limit, page int) (*bookmark.GetBookmarksResult, error) {
	args := m.Called(ctx, userID, limit, page)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookmark.GetBookmarksResult), args.Error(1)
}

func (m *mockBookmarkService) UpdateBookmark(ctx context.Context, bookmarkID, userID string, url, description string) error {
	args := m.Called(ctx, bookmarkID, userID, url, description)
	return args.Error(0)
}

func (m *mockBookmarkService) DeleteBookmark(ctx context.Context, bookmarkID, userID string) error {
	args := m.Called(ctx, bookmarkID, userID)
	return args.Error(0)
}

func TestBookmarkCacheService_GetBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                string
		setupService        func(ctx context.Context) *mockBookmarkService
		setupCache          func(ctx context.Context) *cacheMocks.Repo
		userID              string
		limit               int
		page                int
		expectedBookmarks   int
		expectedCount       int64
		expectedError       bool
	}{
		{
			name: "success - get from cache",
			setupService: func(ctx context.Context) *mockBookmarkService {
				mockSvc := new(mockBookmarkService)
				// Service should NOT be called because data is in cache
				return mockSvc
			},
			setupCache: func(ctx context.Context) *cacheMocks.Repo {
				mockCache := new(cacheMocks.Repo)

				// Cache hit
				cachedData := &bookmark.GetBookmarksResult{
					Bookmarks: []*model.Bookmark{
						{Base: model.Base{ID: "1"}, Url: "https://google.com", Description: "Google"},
						{Base: model.Base{ID: "2"}, Url: "https://github.com", Description: "GitHub"},
					},
					Count: 2,
				}
				dataBytes, _ := json.Marshal(cachedData)

				mockCache.On("GetCacheData", ctx, "get_bookmarks_user_user123", "page_1_limit_10").
					Return(dataBytes, nil)

				return mockCache
			},
			userID:            "user123",
			limit:             10,
			page:              1,
			expectedBookmarks: 2,
			expectedCount:     2,
			expectedError:     false,
		},
		{
			name: "success - cache miss, get from service",
			setupService: func(ctx context.Context) *mockBookmarkService {
				mockSvc := new(mockBookmarkService)

				// Service called on cache miss
				mockSvc.On("GetBookmarks", ctx, "user456", 10, 1).
					Return(&bookmark.GetBookmarksResult{
						Bookmarks: []*model.Bookmark{
							{Base: model.Base{ID: "1"}, Url: "https://golang.org", Description: "Go"},
						},
						Count: 1,
					}, nil)

				return mockSvc
			},
			setupCache: func(ctx context.Context) *cacheMocks.Repo {
				mockCache := new(cacheMocks.Repo)

				// Cache miss
				mockCache.On("GetCacheData", ctx, "get_bookmarks_user_user456", "page_1_limit_10").
					Return(nil, nil) // empty cache

				// Cache set after getting from service
				mockCache.On("SetCacheData", ctx, "get_bookmarks_user_user456", "page_1_limit_10",
					mock.MatchedBy(func(data []byte) bool {
						return len(data) > 0
					}), 24*time.Hour).
					Return(nil)

				return mockCache
			},
			userID:            "user456",
			limit:             10,
			page:              1,
			expectedBookmarks: 1,
			expectedCount:     1,
			expectedError:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			mockSvc := tc.setupService(ctx)
			mockCache := tc.setupCache(ctx)

			// Create cache service wrapper
			cacheService := bookmark.NewServiceWithCache(mockSvc, mockCache)

			result, err := cacheService.GetBookmarks(ctx, tc.userID, tc.limit, tc.page)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedBookmarks, len(result.Bookmarks))
				assert.Equal(t, tc.expectedCount, result.Count)
			}

			mockSvc.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}
