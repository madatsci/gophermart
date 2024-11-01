package app

import (
	"context"

	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/internal/app/logger"
	"github.com/madatsci/gophermart/internal/app/server"
	"github.com/madatsci/gophermart/internal/app/store"
	"go.uber.org/zap"
)

type (
	App struct {
		config *config.Config
		server *server.Server
		store  store.Store
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

	store, err := store.New(ctx, config)
	if err != nil {
		return nil, err
	}

	srv := server.New(config, log)

	app := &App{
		config: config,
		logger: log,
		store:  store,
		server: srv,
	}

	return app, nil
}

// Start starts the application.
func (a *App) Start() error {
	return a.server.Start()
}

// Store is used for migrations.
func (a *App) Store() store.Store {
	return a.store
}
