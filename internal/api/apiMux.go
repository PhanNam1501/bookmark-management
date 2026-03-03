package api

import (
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/handler"
	"github.com/PhanNam1501/bookmark-management/internal/service"
)

type EngineMux interface {
	StartMux() error
}

type apiMux struct {
	mux *http.ServeMux
}

func NewMux() EngineMux {
	a := &apiMux{
		mux: http.NewServeMux(),
	}
	a.registerEPMux()
	return a
}

func (a *apiMux) StartMux() error {
	return http.ListenAndServe(":8080", a.mux)
}

func (a *apiMux) registerEPMux() {
	passSvc := service.NewPassword()
	passHandler := handler.NewPassword(passSvc)
	a.mux.HandleFunc("/gen-pass-mux", passHandler.GenPassForMux)
}
