package server

import (
	"context"
	"net/http"

	"github.com/v7ktory/test/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Cfg, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.Server.Port,
			Handler:        handler,
			ReadTimeout:    cfg.Server.ReadTimeout,
			WriteTimeout:   cfg.Server.WriteTimeout,
			MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
