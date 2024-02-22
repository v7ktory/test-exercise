package rest

import (
	"net/http"

	"github.com/v7ktory/test/internal/service"
)

type Handler struct {
	Svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{
		Svc: svc,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/signup", h.SignUp)
	return mux
}
