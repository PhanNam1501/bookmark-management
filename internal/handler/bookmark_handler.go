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

// GenUuid godoc
// @Summary Health check + generate instance UUID
// @Description Check Redis connection and return a unique instance ID if service is healthy
// @Tags Bookmark
// @Accept json
// @Produce json
// @Success 200 {object} BaseResponse "Service is healthy"
// @Failure 500 {object} map[string]string "Redis connection failed"
// @Router /health-check [get]
func (h *bookmarkHandler) GenUuid(c *gin.Context) {
	err := h.svc.CheckRedisConnection(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	uuid := h.svc.GenerateUuid()
	resp := BaseResponse{
		Message:     "OK",
		ServiceName: "bookmark_service",
		InstanceID:  uuid,
	}

	c.JSON(http.StatusOK, resp)
}
