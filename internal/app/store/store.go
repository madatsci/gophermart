package store

import (
	"context"
	"errors"

	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/internal/app/database"
	"github.com/madatsci/gophermart/internal/app/models"
	db "github.com/madatsci/gophermart/internal/app/store/database"
)

type Store interface {
	// Users
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)

	// Accounts
	CreateAccount(ctx context.Context, account models.Account) (models.Account, error)
	GetAccountByUserID(ctx context.Context, userID string) (models.Account, error)

	// Orders
	CreateOrder(ctx context.Context, order models.Order) (models.Order, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (models.Order, error)
}

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
