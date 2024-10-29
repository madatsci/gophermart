package app

import (
	"context"

	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/internal/app/logger"
	"github.com/madatsci/gophermart/internal/app/server"
	"go.uber.org/zap"
)

type (
	App struct {
		config *config.Config
		server *server.Server
		logger *zap.SugaredLogger
	}

	Options struct {
		RunAddress           string
		AccrualSystemAddress string
		DatabaseURI          string
	}
)

// New creates new App.
func New(ctx context.Context, opts Options) (*App, error) {
	config := config.New(opts.RunAddress, opts.AccrualSystemAddress, opts.DatabaseURI)

	log, err := logger.New()
	if err != nil {
		return nil, err
	}

	srv := server.New(config, log)

	app := &App{
		config: config,
		logger: log,
		server: srv,
	}

	return app, nil
}

// Start starts the application.
func (a *App) Start() error {
	return a.server.Start()
}
