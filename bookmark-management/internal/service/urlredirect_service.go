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
}

func NewUrlRedirect(urlStorage repository.URLStorage) URLRedirect {
	return &urlRedirectService{
		urlStorage: urlStorage,
	}
}

func (u *urlRedirectService) GetRedirectURL(ctx context.Context, code string) (string, error) {
	val, err := u.urlStorage.GetRedirectURL(ctx, code)
	return val, err
}
