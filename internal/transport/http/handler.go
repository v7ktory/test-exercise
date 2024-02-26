package http

import (
	"net/http"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()

	r.HandleFunc("/auth/signup", h.SignUp).Methods("POST")
	r.HandleFunc("/auth/login", h.Login).Methods("POST")
	r.HandleFunc("/auth/refresh", h.Refresh).Methods("POST")

	return r
}
