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
	mux.HandleFunc("/auth/login", h.Login)
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
