package store

import (
	"context"
	"errors"

	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/internal/app/database"
	db "github.com/madatsci/gophermart/internal/app/store/database"
)

type Store interface{}

func New(ctx context.Context, cfg *config.Config) (Store, error) {
	if cfg.DatabaseURI != "" {
		conn, err := database.NewClient(ctx, cfg.DatabaseURI)
		if err != nil {
			return nil, err
		}
		return db.New(ctx, conn)
	}

	return nil, errors.New("database URI must be provided")
}
