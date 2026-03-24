package service

import (
	"github.com/google/uuid"
)

//go:generate mockery --name Bookmark --filename bookmark_service.go
type bookmarkService struct {
}

type Bookmark interface {
	GenerateUuid() string
}

func NewBookmark() Bookmark {
	return &bookmarkService{}
}

func (b *bookmarkService) GenerateUuid() string {
	return uuid.New().String()
}
