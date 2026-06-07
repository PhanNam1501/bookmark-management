package bookmark

import (
	"github.com/PhanNam1501/bookmark-management/internal/service/bookmark"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	CreateBookmark(c *gin.Context)
	GetBookmarks(c *gin.Context)
	UpdateBookmark(c *gin.Context)
	DeleteBookmark(c *gin.Context)
}

type handler struct {
	s bookmark.Service
}

func NewHandler(s bookmark.Service) Handler {
	return &handler{
		s: s,
	}
}
