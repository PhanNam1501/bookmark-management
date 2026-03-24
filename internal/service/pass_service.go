package service

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789"
	passLen = 10
)

type passwordService struct {
}

// Password interface represents the password service
//
//go:generate mockery --name Password --filename pass_service.go
type Password interface {
	GeneratePassword() (string, error)
}

// Return a new instance of the password service
func NewPassword() Password {
	return &passwordService{}
}

// Generate a secure random password
func (s *passwordService) GeneratePassword() (string, error) {
	var strBuilder bytes.Buffer

	for i := 0; i < passLen; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		strBuilder.WriteByte(charset[randomIndex.Int64()])
	}

	return strBuilder.String(), nil
}
