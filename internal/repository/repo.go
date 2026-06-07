package repository

//go:generate mockery --name Repository --filename repository_mock.go

import (
	"context"

	"github.com/PhanNam1501/bookmark-management/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	CreateBookmark(ctx context.Context, b *model.Bookmark) error
	QueryBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, error)
	GetBookmarkCounts(ctx context.Context, userID string) (int64, error)
	UpdateBookmark(ctx context.Context, bookmarkID, userID string, updateData map[string]interface{}) error
	DeleteBookmark(ctx context.Context, bookmarkID, userID string) error
	GetBookmarkByCode(ctx context.Context, code string) (*model.Bookmark, error)
	UpdateBookmarkCode(ctx context.Context, bookmarkID, code string) error
}

type bookmarkRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &bookmarkRepo{
		db: db,
	}
}
