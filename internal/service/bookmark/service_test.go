package bookmark

import (
	"context"
	"testing"

	"github.com/PhanNam1501/bookmark-management/internal/model"
	repomocks "github.com/PhanNam1501/bookmark-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockKeyGen struct {
	mock.Mock
}

func (m *mockKeyGen) Generate() string {
	args := m.Called()
	return args.String(0)
}

func TestGetBookmarks(t *testing.T) {
	mockRepo := new(repomocks.Repository)
	mockKeyGen := new(mockKeyGen)

	ctx := context.Background()
	userID := "bb4b022a-03b5-4d9a-8edf-a19e79f8fd60"
	limit := 10
	page := 1
	offset := (page - 1) * limit // offset = 0 for page 1

	// Expected bookmarks
	mockBookmarks := []*model.Bookmark{
		{
			Base: model.Base{ID: "1"},
			Url:  "https://google.com",
			Description: "Google",
			UserId: userID,
		},
		{
			Base: model.Base{ID: "2"},
			Url:  "https://github.com",
			Description: "GitHub",
			UserId: userID,
		},
		{
			Base: model.Base{ID: "3"},
			Url:  "https://golang.org",
			Description: "Go",
			UserId: userID,
		},
	}

	// Setup expectations - verify offset is calculated correctly
	mockRepo.On("QueryBookmarks", ctx, userID, limit, offset).Return(mockBookmarks, nil)
	mockRepo.On("GetBookmarkCounts", ctx, userID).Return(int64(3), nil)

	// Create service with mock
	service := NewService(mockRepo, mockKeyGen)
	result, err := service.GetBookmarks(ctx, userID, limit, page)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result.Bookmarks))
	assert.Equal(t, int64(3), result.Count)
	assert.Equal(t, "Google", result.Bookmarks[0].Description)
	assert.Equal(t, "GitHub", result.Bookmarks[1].Description)

	// Verify mock was called with correct parameters
	mockRepo.AssertExpectations(t)
}

func TestGetBookmarksPage2(t *testing.T) {
	mockRepo := new(repomocks.Repository)
	mockKeyGen := new(mockKeyGen)

	ctx := context.Background()
	userID := "test-user"
	limit := 10
	page := 2
	offset := (page - 1) * limit // offset = 10 for page 2

	mockBookmarks := []*model.Bookmark{
		{
			Base: model.Base{ID: "11"},
			Url:  "https://example1.com",
			UserId: userID,
		},
		{
			Base: model.Base{ID: "12"},
			Url:  "https://example2.com",
			UserId: userID,
		},
	}

	// Verify offset calculation is correct: (page - 1) * limit = (2 - 1) * 10 = 10
	mockRepo.On("QueryBookmarks", ctx, userID, limit, offset).Return(mockBookmarks, nil)
	mockRepo.On("GetBookmarkCounts", ctx, userID).Return(int64(12), nil)

	service := NewService(mockRepo, mockKeyGen)
	result, err := service.GetBookmarks(ctx, userID, limit, page)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result.Bookmarks))
	assert.Equal(t, int64(12), result.Count)

	// Verify QueryBookmarks was called with correct offset (10)
	mockRepo.AssertCalled(t, "QueryBookmarks", ctx, userID, limit, offset)
	mockRepo.AssertExpectations(t)
}
