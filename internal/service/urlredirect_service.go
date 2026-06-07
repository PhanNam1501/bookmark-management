package service

import (
	"context"

	"github.com/PhanNam1501/bookmark-management/internal/repository"
)

//go:generate mockery --name URLRedirect --filename urlredirect_service.go
type URLRedirect interface {
	GetRedirectURL(ctx context.Context, code string) (string, error)
}

type urlRedirectService struct {
	urlStorage repository.URLStorage
	repo       repository.Repository
}

func NewUrlRedirect(urlStorage repository.URLStorage) URLRedirect {
	return &urlRedirectService{
		urlStorage: urlStorage,
	}
}

func NewUrlRedirectWithDB(urlStorage repository.URLStorage, repo repository.Repository) URLRedirect {
	return &urlRedirectService{
		urlStorage: urlStorage,
		repo:       repo,
	}
}

func (u *urlRedirectService) GetRedirectURL(ctx context.Context, code string) (string, error) {
	// Code length 7 = Redis (URLStorage)
	if len(code) == 7 {
		val, err := u.urlStorage.GetRedirectURL(ctx, code)
		return val, err
	}

	// Code length 8 = Database (Bookmark)
	if len(code) == 8 && u.repo != nil {
		bookmark, err := u.repo.GetBookmarkByCode(ctx, code)
		if err != nil {
			return "", err
		}
		return bookmark.Url, nil
	}

	// Fallback to Redis
	val, err := u.urlStorage.GetRedirectURL(ctx, code)
	return val, err
}
