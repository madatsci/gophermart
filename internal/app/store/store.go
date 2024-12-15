package store

import (
	"context"
	"errors"

	"github.com/madatsci/gophermart/internal/app/models"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Store interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)

	CreateAccount(ctx context.Context, account models.Account) (models.Account, error)
	GetAccountByUserID(ctx context.Context, userID string) (models.Account, error)
	WithdrawBalance(ctx context.Context, userID string, orderNumber string, sum float32) (models.Account, error)
	AddBalance(ctx context.Context, order models.Order) (models.Account, error)

	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByNumber(ctx context.Context, orderNumber string) (models.Order, error)
	ListOrdersByAccountID(ctx context.Context, accountID string, limit int) ([]models.Order, error)
	ListOrdersByStatus(ctx context.Context, statuses []models.OrderStatus, limit int) ([]models.Order, error)
	UpdateOrder(ctx context.Context, order models.Order, prevStatus models.OrderStatus) (models.Order, error)

	GetWithdrawals(ctx context.Context, accountID string, direction models.TxDirection, limit int) ([]models.Transaction, error)
}

type NotEnoughBalanceError struct {
	Err               error
	Balance           float32
	WithdrawRequested float32
}

func (e *NotEnoughBalanceError) Error() string {
	return e.Err.Error()
}

type InsertError struct {
	Err error
}

func (e *InsertError) Error() string {
	return e.Err.Error()
}

func (e *InsertError) IntegrityViolation() bool {
	var pgErr pgdriver.Error
	return errors.As(e.Err, &pgErr) && pgErr.IntegrityViolation()
}

type StoreError interface {
	IntegrityViolation() bool
}
