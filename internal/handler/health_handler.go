package handler

import (
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type healthCheckHandler struct {
	svc service.Id
}

type HealthCheck interface {
	GenId(c *gin.Context)
}

func NewHealthCheck(svc service.Id) HealthCheck {
	return &healthCheckHandler{svc: svc}
}

func (h *healthCheckHandler) GenId(c *gin.Context) {
	id := h.svc.GetId()

	c.String(http.StatusOK, id)
}
