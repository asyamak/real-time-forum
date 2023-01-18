package server

import (
	"context"
	"net/http"
	"real-time-forum/internal/config"
	"time"

	"github.com/rshezarr/gorr"
)

type Server struct {
	server *http.Server
}

func NewServer(cfg *config.Config, router *gorr.Router) *Server {
	return &Server{
		server: &http.Server{
			Addr:           cfg.API.Port,
			Handler:        router,
			MaxHeaderBytes: cfg.API.MaxHeaderBytes << 20,
			ReadTimeout:    time.Duration(cfg.API.ReadTimeout) * time.Second,
			WriteTimeout:   time.Duration(cfg.API.WriteTimeout) * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
