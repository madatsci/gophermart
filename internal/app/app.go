package app

import (
	"context"
	"errors"
	"time"

	"github.com/madatsci/gophermart/internal/app/accrual"
	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/internal/app/database"
	"github.com/madatsci/gophermart/internal/app/logger"
	"github.com/madatsci/gophermart/internal/app/server"
	"github.com/madatsci/gophermart/internal/app/store"
	db "github.com/madatsci/gophermart/internal/app/store/database"
	"go.uber.org/zap"
)

type (
	App struct {
		config *config.Config
		server *server.Server
		store  store.Store
		// TODO use interface
		as     *accrual.AccrualService
		logger *zap.SugaredLogger
	}

	Options struct {
		RunAddress           string
		AccrualSystemAddress string
		DatabaseURI          string
		TokenSecret          []byte
		TokenDuration        time.Duration
	}
)

// New creates new App.
func New(ctx context.Context, opts Options) (*App, error) {
	config := config.New(opts.RunAddress, opts.AccrualSystemAddress, opts.DatabaseURI, opts.TokenSecret, opts.TokenDuration)

	log, err := logger.New()
	if err != nil {
		return nil, err
	}

	store, err := newStore(ctx, config)
	if err != nil {
		return nil, err
	}

	srv := server.New(config, store, log)

	app := &App{
		config: config,
		logger: log,
		store:  store,
		as:     accrual.New(config, store, log),
		server: srv,
	}

	return app, nil
}

// Start starts the application.
func (a *App) Start(ctx context.Context) error {
	go a.syncOrders(ctx)
	return a.server.Start()
}

// Store is used for migrations.
func (a *App) Store() store.Store {
	return a.store
}

func (a *App) syncOrders(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := a.as.SyncOrders(ctx)
			if err != nil {
				a.logger.With("err", err).Errorln("could not sync orders")
			}
		}
	}
}

func newStore(ctx context.Context, cfg *config.Config) (store.Store, error) {
	if cfg.DatabaseURI != "" {
		conn, err := database.NewClient(ctx, cfg.DatabaseURI)
		if err != nil {
			return nil, err
		}
		return db.New(ctx, conn)
	}

	return nil, errors.New("database URI must be provided")
}
