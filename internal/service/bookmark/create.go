package bookmark

import (
	"context"

	"github.com/PhanNam1501/bookmark-management/internal/model"
	"github.com/PhanNam1501/bookmark-management/pkg/utils"
)

func (s *bookmarkService) CreateBookmark(ctx context.Context, userId, url, description string) (*model.Bookmark, error) {
	b := &model.Bookmark{
		Description: description,
		Url:         url,
		Code:        "", // Will be set after insert
		UserId:      userId,
	}

	err := s.r.CreateBookmark(ctx, b)
	if err != nil {
		return nil, err
	}

	// Generate code from code_int (auto_increment from DB)
	code := utils.GenerateBookmarkCode(b.CodeInt)
	code = utils.MapGenerateCodeForBookmark(code)

	// Update bookmark with generated code
	err = s.r.UpdateBookmarkCode(ctx, b.ID, code)
	if err != nil {
		return nil, err
	}

	b.Code = code
	return b, nil
}

func (s *bookmarkService) UpdateBookmark(ctx context.Context, bookmarkID, userID string, url, description string) error {
	updateData := map[string]interface{}{
		"url":         url,
		"description": description,
	}

	return s.r.UpdateBookmark(ctx, bookmarkID, userID, updateData)
}

func (s *bookmarkService) DeleteBookmark(ctx context.Context, bookmarkID, userID string) error {
	return s.r.DeleteBookmark(ctx, bookmarkID, userID)
}
