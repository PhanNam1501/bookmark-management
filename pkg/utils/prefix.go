package utils

import (
	"math/rand"

	"github.com/PhanNam1501/bookmark-management/pkg/algorithm"
)

var redisPrefixes = []byte("abcdefg")
var sqlPrefixes = []byte("hijklmnopqrstuvwxyz")

func GenerateBookmarkCode(codeInt int) string {
	code := algorithm.CreateCode(codeInt)
	return code
}

func GetPrefixFromCode(code string) string {
	return string(code[0])
}

func MapGenerateCodeForRedis(code string) string {
	prefix := redisPrefixes[rand.Intn(len(redisPrefixes))]
	return string(prefix) + code
}

func MapGenerateCodeForBookmark(code string) string {
	prefix := sqlPrefixes[rand.Intn(len(sqlPrefixes))]
	return string(prefix) + code
}
