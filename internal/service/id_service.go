package service

import "github.com/google/uuid"

type idService struct {
}

//go:generate mockery --name Id --filename id_service.go
type Id interface {
	GetId() string
}

func NewId() Id {
	return &idService{}
}

func (i *idService) GetId() string {
	id := uuid.New().String()
	return id
}
