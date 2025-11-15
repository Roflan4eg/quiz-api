package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Roflan4eg/quiz-api/config"
	"github.com/Roflan4eg/quiz-api/pkg/logger"
)

type Server struct {
	server *http.Server
	name   string
}

func New(handler http.Handler, cfg *config.Config) *Server {
	return &Server{
		server: &http.Server{
			Addr:         ":" + cfg.HTTP.Port,
			Handler:      handler,
			ReadTimeout:  cfg.HTTP.ReadTimeout * time.Second,
			WriteTimeout: cfg.HTTP.WriteTimeout * time.Second,
		},
		name: "quiz-api",
	}
}

func (s *Server) Start() error {
	logger.Info("starting server",
		"name", s.name,
		"address", s.server.Addr,
	)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Name() string {
	return s.name
}
