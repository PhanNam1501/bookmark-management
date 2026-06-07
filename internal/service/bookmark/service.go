package bookmark

//go:generate mockery --name Service --filename service_mock.go

import (
	"context"

	"github.com/PhanNam1501/bookmark-management/internal/model"
	"github.com/PhanNam1501/bookmark-management/internal/repository"
	"github.com/PhanNam1501/bookmark-management/pkg/utils"
)

type Service interface {
	CreateBookmark(ctx context.Context, userid, url, description string) (*model.Bookmark, error)
	GetBookmarks(ctx context.Context, userID string, limit, page int) (*GetBookmarksResult, error)
	UpdateBookmark(ctx context.Context, bookmarkID, userID string, url, description string) error
	DeleteBookmark(ctx context.Context, bookmarkID, userID string) error
}

type bookmarkService struct {
	r      repository.Repository
	keyGen utils.KeyGenerator
}

func NewService(r repository.Repository, keyGen utils.KeyGenerator) Service {
	return &bookmarkService{
		r:      r,
		keyGen: keyGen,
	}
}
