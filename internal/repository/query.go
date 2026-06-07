package repository

import (
	"context"

	"github.com/PhanNam1501/bookmark-management/internal/model"
	"github.com/PhanNam1501/bookmark-management/pkg/dbutils"
)

func (r *bookmarkRepo) QueryBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, error) {
	var bookmarks []*model.Bookmark
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&bookmarks).Error; err != nil {
		return nil, err
	}

	return bookmarks, nil
}

func (r *bookmarkRepo) GetBookmarkCounts(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Bookmark{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, dbutils.CatchDBError(err)
	}
	return count, nil
}

func (r *bookmarkRepo) GetBookmarkByCode(ctx context.Context, code string) (*model.Bookmark, error) {
	var bookmark *model.Bookmark
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&bookmark).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return bookmark, nil
}
