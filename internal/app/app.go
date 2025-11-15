package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/Roflan4eg/quiz-api/config"
	"github.com/Roflan4eg/quiz-api/internal/app/http"
	"github.com/Roflan4eg/quiz-api/internal/handler"
	"github.com/Roflan4eg/quiz-api/internal/repository"
	"github.com/Roflan4eg/quiz-api/internal/service"
	"github.com/Roflan4eg/quiz-api/internal/storage"
	"github.com/Roflan4eg/quiz-api/pkg/logger"
)

type Server interface {
	Start() error
	Stop(ctx context.Context) error
	Name() string
}

type App struct {
	cfg        *config.Config
	storage    *storage.Container
	repository *repository.Container
	services   *service.Container
	handlers   *handler.Handler
	servers    []Server
	closer     *Closer
}

func New(cfg *config.Config) *App {
	return &App{
		cfg:    cfg,
		closer: &Closer{},
	}
}

func (a *App) Setup() error {
	stor, err := storage.NewContainer(a.cfg)
	if err != nil {
		return fmt.Errorf("storage setup: %w", err)
	}
	a.storage = stor
	a.closer.Add(a.storage.Close)

	a.repository = repository.NewContainer(a.storage)
	a.services = service.NewContainer(a.repository)
	a.handlers = handler.NewHandler(a.services)

	if err = a.setupServers(); err != nil {
		return fmt.Errorf("server setup: %w", err)
	}

	a.closer.Add(a.stopServers)

	return nil
}

func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, server := range a.servers {
		go func(server Server) {
			if err := server.Start(); err != nil {
				logger.Error("Failed to start server",
					"server", server.Name(),
					"error", err.Error(),
				)
			}
		}(server)
	}
	go StartShutdownListener(ctx, cancel)

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.App.ShutdownTimeout)
	defer cancel()

	if err := a.closer.Close(shutdownCtx); err != nil {
		return err
	}

	return nil
}

func (a *App) setupServers() error {
	a.servers = append(a.servers, http.New(a.handlers.InitRoutes(), a.cfg))
	return nil
}

func (a *App) stopServers(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(len(a.servers))
	for _, server := range a.servers {
		go func(server Server) {
			logger.Info("Stopping server", "server", server.Name())
			if err := server.Stop(ctx); err != nil {
				logger.Warn("Failed to stop server",
					"server", server.Name(),
					"error", err,
				)
			}
			wg.Done()
		}(server)
	}
	wg.Wait()
	return nil
}
