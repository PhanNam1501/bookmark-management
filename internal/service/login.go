package service

import (
	"context"
	"errors"
	"time"

	"github.com/PhanNam1501/bookmark-management/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

const tokenExp = time.Hour * 24

var (
	errClientErr = errors.New("wrong username or password")
)

func (u *userService) Login(ctx context.Context, username, password string) (string, error) {
	// Get user from repo, check if user exist
	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	// check if password match
	passwordMatched := utils.VerifyPassword(user.Password, password)
	if !passwordMatched {
		return "", errClientErr
	}
	// create token
	jwtContent := jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(tokenExp).Unix(),
	}
	// return token
	jwtString, err := u.jwtGen.GenerateToken(jwtContent)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}
