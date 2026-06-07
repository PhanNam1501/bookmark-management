package bookmark

import (
	"context"

	"github.com/PhanNam1501/bookmark-management/internal/model"
)

type GetBookmarksResult struct {
	Bookmarks []*model.Bookmark
	Count     int64
}

func (s *bookmarkService) GetBookmarks(ctx context.Context, userID string, limit, page int) (*GetBookmarksResult, error) {
	offset := (page - 1) * limit

	bookmarks, err := s.r.QueryBookmarks(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	count, err := s.r.GetBookmarkCounts(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &GetBookmarksResult{
		Bookmarks: bookmarks,
		Count:     count,
	}, nil
}
