package server

import (
	"context"
	"net/http"
	"real-time-forum/internal/config"
	"time"

	"github.com/rshezarr/gorr"
)

type Server struct {
	server            *http.Server
	shutdownTimeout   time.Duration
	ServerErrorNotify chan error
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
		shutdownTimeout:   time.Duration(cfg.API.ShutdownTimeout) * time.Second,
		ServerErrorNotify: make(chan error, 1),
	}
}

func (s *Server) Start() {
	s.ServerErrorNotify <- s.server.ListenAndServe()
}

func (s *Server) ServerErrNotify() <-chan error {
	return s.ServerErrorNotify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
