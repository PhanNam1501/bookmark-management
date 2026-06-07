package utils

import (
	"crypto/rand"
	"encoding/hex"
)

type KeyGenerator interface {
	Generate() string
}

type randomKeyGen struct {
	length int
}

func NewRandomKeyGenerator(length int) KeyGenerator {
	return &randomKeyGen{length: length}
}

func (rk *randomKeyGen) Generate() string {
	bytes := make([]byte, rk.length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
