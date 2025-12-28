package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	port   int
	server *http.Server
}

func NewServer(logger *zap.Logger, port int, handler http.Handler) *Server {
	return &Server{
		logger: logger,
		port:   port,
		server: &http.Server{
			Addr:              fmt.Sprintf(":%d", port),
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      20 * time.Second,
			IdleTimeout:       120 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	s.logger.Info("starting server", zap.Int("port", s.port))
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")
	return s.server.Shutdown(ctx)
}
