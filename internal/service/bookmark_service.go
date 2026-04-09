package service

import (
	"context"

	"github.com/PhanNam1501/bookmark-management/internal/repository"
	"github.com/google/uuid"
)

//go:generate mockery --name Bookmark --filename bookmark_service.go
type bookmarkService struct {
	urlStorage repository.URLStorage
}

type Bookmark interface {
	GenerateUuid() string
	CheckRedisConnection(ctx context.Context) error
}

func NewBookmark(urlStorage repository.URLStorage) Bookmark {
	return &bookmarkService{
		urlStorage: urlStorage,
	}
}

func (b *bookmarkService) GenerateUuid() string {
	return uuid.New().String()
}

func (b *bookmarkService) CheckRedisConnection(ctx context.Context) error {
	err := b.urlStorage.CheckRedisConnection(ctx)
	if err != nil {
		return err
	}

	return nil

}
