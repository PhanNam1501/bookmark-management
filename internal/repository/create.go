package repository

import (
	"context"

	"github.com/PhanNam1501/bookmark-management/internal/model"
	"github.com/PhanNam1501/bookmark-management/pkg/dbutils"
)

func (r *bookmarkRepo) CreateBookmark(ctx context.Context, b *model.Bookmark) error {
	err := r.db.WithContext(ctx).Create(b).Error
	if err != nil {
		return dbutils.CatchDBError(err)
	}

	return nil
}

func (r *bookmarkRepo) UpdateBookmark(ctx context.Context, bookmarkID, userID string, updateData map[string]interface{}) error {
	err := r.db.WithContext(ctx).
		Model(&model.Bookmark{}).
		Where("id = ? AND user_id = ?", bookmarkID, userID).
		Updates(updateData).Error

	if err != nil {
		return dbutils.CatchDBError(err)
	}

	return nil
}

func (r *bookmarkRepo) DeleteBookmark(ctx context.Context, bookmarkID, userID string) error {
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", bookmarkID, userID).
		Delete(&model.Bookmark{}).Error

	if err != nil {
		return dbutils.CatchDBError(err)
	}

	return nil
}

func (r *bookmarkRepo) UpdateBookmarkCode(ctx context.Context, bookmarkID, code string) error {
	err := r.db.WithContext(ctx).
		Model(&model.Bookmark{}).
		Where("id = ?", bookmarkID).
		Update("code", code).Error

	if err != nil {
		return dbutils.CatchDBError(err)
	}

	return nil
}
