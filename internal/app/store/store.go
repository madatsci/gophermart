package store

import (
	"context"

	"github.com/madatsci/gophermart/internal/app/models"
	"github.com/shopspring/decimal"
)

type Store interface {
	// Users
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)

	// Accounts
	CreateAccount(ctx context.Context, account models.Account) (models.Account, error)
	GetAccountByUserID(ctx context.Context, userID string) (models.Account, error)
	WithdrawBalance(ctx context.Context, userID string, orderNumber string, sum decimal.Decimal) (models.Account, error)

	// Orders
	CreateOrder(ctx context.Context, order models.Order) (models.Order, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (models.Order, error)
	ListOrdersByAccountID(ctx context.Context, accountID string, limit int) ([]models.Order, error)
	ListOrdersByStatus(ctx context.Context, statuses []models.OrderStatus, limit int) ([]models.Order, error)
	UpdateOrder(ctx context.Context, order models.Order, prevStatus models.OrderStatus) (models.Order, error)

	// Transactions
	GetWithdrawals(ctx context.Context, accountID string, direction models.TxDirection, limit int) ([]models.Transaction, error)
}

type NotEnoughBalanceError struct {
	Err               error
	Balance           decimal.Decimal
	WithdrawRequested decimal.Decimal
}

func (e *NotEnoughBalanceError) Error() string {
	return e.Err.Error()
}
