package handler

import (
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}

type bookmarkHandler struct {
	svc service.Bookmark
}

type Bookmark interface {
	GenUuid(c *gin.Context)
}

func NewBookmarkHandler(svc service.Bookmark) Bookmark {
	return &bookmarkHandler{svc: svc}
}

func (h *bookmarkHandler) GenUuid(c *gin.Context) {
	uuid := h.svc.GenerateUuid()
	resp := BaseResponse{
		Message:     "OK",
		ServiceName: "bookmark_service",
		InstanceID:  uuid,
	}
	c.JSON(http.StatusOK, resp)
}
