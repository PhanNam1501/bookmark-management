package handler

import (
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type passwordHandler struct {
	svc service.Password
}

type Password interface {
	GenPass(c *gin.Context)
}

func NewPasswordHandler(svc service.Password) Password {
	return &passwordHandler{svc: svc}
}

// @Summary Generate password
// @Tags Password
// @Produce json
// @Success 200 {object} string
// @Router /gen-pass [get]
func (h *passwordHandler) GenPass(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		c.String(http.StatusInternalServerError, "err")
		return
	}

	c.String(http.StatusOK, pass)
}
