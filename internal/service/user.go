package service

import (
	"context"

	"github.com/PhanNam1501/bookmark-management/internal/model"
	"github.com/PhanNam1501/bookmark-management/internal/repository"
	"github.com/PhanNam1501/bookmark-management/pkg/jwtutils"
	"github.com/PhanNam1501/bookmark-management/pkg/utils"
)

type User interface {
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)
	Login(ctx context.Context, username, password string) (string, error)
	GetCurrentUser(ctx context.Context, id string) (*model.User, error)
}

type userService struct {
	repo   repository.User
	jwtGen jwtutils.JWTGenerator
}

func NewUser(repo repository.User, jwtGen jwtutils.JWTGenerator) User {
	return &userService{
		repo:   repo,
		jwtGen: jwtGen,
	}
}

func (u *userService) CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error) {
	hashPwd := utils.HashPassword(password)

	newuser := model.User{
		UserName:    username,
		DisplayName: displayName,
		Email:       email,
		Password:    hashPwd,
	}

	res, err := u.repo.CreateUser(ctx, &newuser)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *userService) GetCurrentUser(ctx context.Context, id string) (*model.User, error) {
	return u.repo.GetUserByID(ctx, id)
}
