package handler

import (
	"fmt"
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type passwordHandler struct {
	svc service.Password
}

type Password interface {
	GenPass(c *gin.Context)
	GenPassForMux(w http.ResponseWriter, r *http.Request)
}

func NewPassword(svc service.Password) Password {
	return &passwordHandler{svc: svc}
}

func (h *passwordHandler) GenPass(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		c.String(http.StatusInternalServerError, "err")
	}

	c.String(http.StatusOK, pass)
}

func (h *passwordHandler) GenPassForMux(w http.ResponseWriter, r *http.Request) {
	fmt.Println("for mux")

	pass, err := h.svc.GeneratePassword()
	if err != nil {
		http.Error(w, "cannot generate password", http.StatusNotFound)
	}
	_, err = w.Write([]byte(pass))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
