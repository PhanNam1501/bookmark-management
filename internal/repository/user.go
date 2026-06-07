package repository

import (
	"context"

	"github.com/PhanNam1501/bookmark-management/internal/model"
	"github.com/PhanNam1501/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

type User interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

type user struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) User {
	return &user{
		db: db,
	}
}

func (u *user) CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) {
	err := u.db.WithContext(ctx).Create(newUser).Error

	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return newUser, nil
}

func (u *user) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}

	err := u.db.WithContext(ctx).Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return user, nil
}

func (u *user) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}

	err := u.db.WithContext(ctx).Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return user, nil
}
