package bookmark

import (
	"github.com/PhanNam1501/bookmark-management/internal/service/bookmark"
	"github.com/PhanNam1501/bookmark-management/internal/service/queue"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	CreateBookmark(c *gin.Context)
	GetBookmarks(c *gin.Context)
	UpdateBookmark(c *gin.Context)
	DeleteBookmark(c *gin.Context)
	Import(c *gin.Context)
}

type handler struct {
	s bookmark.Service
	q queue.Service
}

func NewHandler(s bookmark.Service, q queue.Service) Handler {
	return &handler{
		s: s,
		q: q,
	}
}
