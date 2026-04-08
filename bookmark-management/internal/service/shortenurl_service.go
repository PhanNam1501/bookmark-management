package service

import (
	"context"
	"fmt"

	"github.com/PhanNam1501/bookmark-management/internal/repository"
)

const (
	retryAttemps = 5
)

//go:generate mockery --name ShortenURL --filename shortenurl_service.go
type ShortenURL interface {
	ShortenURL(ctx context.Context, url string) (string, error)
	LinkShortenURL(ctx context.Context, url string, exp int) (string, error)
}

type shortenURLService struct {
	urlStorage repository.URLStorage
	codegen    Password
}

func NewShortenURL(urlStorage repository.URLStorage, codegen Password) ShortenURL {
	return &shortenURLService{
		urlStorage: urlStorage,
		codegen:    codegen,
	}
}

func (s *shortenURLService) ShortenURL(ctx context.Context, url string) (string, error) {
	//gen code
	for i := 0; i < retryAttemps; i++ {
		code, err := s.codegen.GeneratePassword()
		if err != nil {
			return "", err
		}
		ok, err := s.urlStorage.StoreURL(ctx, code, url)
		if err != nil {
			return "", fmt.Errorf("StoreURL: %w", err)
		}
		if ok == "OK" {
			return code, nil
		}
	}

	return "", fmt.Errorf("Cannot generate code after %d", retryAttemps)
}

func (s *shortenURLService) LinkShortenURL(ctx context.Context, url string, exp int) (string, error) {
	for i := 0; i < retryAttemps; i++ {
		code, err := s.codegen.GeneratePassword()
		if err != nil {
			return "", err
		}
		ok, err := s.urlStorage.LinkShortenURL(ctx, code, url, exp)
		if err != nil {
			return "", fmt.Errorf("StoreURL: %w", err)
		}
		if ok == "OK" {
			return code, nil
		}
	}

	return "", fmt.Errorf("Cannot generate code after %d", retryAttemps)
}
