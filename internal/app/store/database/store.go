package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/madatsci/gophermart/internal/app/models"
	"github.com/madatsci/gophermart/internal/app/store"
	"github.com/uptrace/bun"
)

type Store struct {
	conn *bun.DB
}

// New creates a new database-driven storage.
func New(ctx context.Context, conn *bun.DB) (*Store, error) {
	store := &Store{conn: conn}
	if err := store.bootstrap(ctx); err != nil {
		return nil, err
	}

	return store, nil
}

// CreateUser saves new user in database.
func (s *Store) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	var result models.User

	err := s.conn.NewInsert().Model(&user).Returning("*").Scan(ctx, &result)

	return result, err
}

// GetUserByLogin fetches user from database by login.
func (s *Store) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	var result models.User

	err := s.conn.NewSelect().Model(&result).Where("login = ?", login).Scan(ctx)

	return result, err
}

// CreateAccount creates new account.
func (s *Store) CreateAccount(ctx context.Context, account models.Account) (models.Account, error) {
	var result models.Account

	err := s.conn.NewInsert().Model(&account).Returning("*").Scan(ctx, &result)

	return result, err
}

// GetAccountByUserID fetches user account by user ID.
func (s *Store) GetAccountByUserID(ctx context.Context, userID string) (models.Account, error) {
	var result models.Account

	err := s.conn.NewSelect().Model(&result).Where("user_id = ?", userID).Scan(ctx)

	return result, err
}

// CreateOrder saves new order in database.
func (s *Store) CreateOrder(ctx context.Context, order *models.Order) error {
	return s.conn.NewInsert().Model(&order).Returning("*").Scan(ctx, &order)
}

// GetOrderByNumber fetches order by its number.
func (s *Store) GetOrderByNumber(ctx context.Context, orderNumber string) (models.Order, error) {
	var result models.Order

	err := s.conn.NewSelect().
		Model(&result).
		Where("number = ?", orderNumber).
		Relation("Account").
		Scan(ctx)

	return result, err
}

// ListOrdersByAccountID fetches orders linked to the account.
func (s *Store) ListOrdersByAccountID(ctx context.Context, accountID string, limit int) ([]models.Order, error) {
	var result []models.Order

	err := s.conn.NewSelect().
		Model(&result).
		Where("account_id = ?", accountID).
		Order("created_at DESC").
		Limit(limit).
		Scan(ctx)

	return result, err
}

// ListOrdersByStatus fetches orders in specified statuses.
func (s *Store) ListOrdersByStatus(ctx context.Context, statuses []models.OrderStatus, limit int) ([]models.Order, error) {
	var result []models.Order

	err := s.conn.NewSelect().
		Model(&result).
		Where("status IN (?)", bun.In(statuses)).
		Order("created_at ASC").
		Limit(limit).
		Scan(ctx)

	return result, err
}

// UpdateOrder updates order in database.
func (s *Store) UpdateOrder(ctx context.Context, order models.Order, prevStatus models.OrderStatus) (models.Order, error) {
	var checkOrder models.Order

	tx, err := s.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return checkOrder, err
	}

	err = tx.NewSelect().
		Model(&checkOrder).
		Where("id = ?", order.ID).
		For("UPDATE").
		Scan(ctx)
	if err != nil {
		tx.Rollback() //nolint:errcheck
		return checkOrder, err
	}
	if checkOrder.Status != prevStatus {
		tx.Rollback() //nolint:errcheck
		return checkOrder, errors.New("sql: update conflict")
	}

	_, err = tx.NewUpdate().
		Model(&order).
		WherePK().
		Column("status", "accrual", "updated_at").
		Returning("*").
		Exec(ctx)
	if err != nil {
		tx.Rollback() //nolint:errcheck
		return checkOrder, err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback() //nolint:errcheck
		return checkOrder, err
	}

	return order, nil
}

// WithdrawBalance withdraws points from balance if there are enough points.
func (s *Store) WithdrawBalance(ctx context.Context, userID string, orderNumber string, sum float32) (models.Account, error) {
	var acc models.Account

	tx, err := s.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return acc, err
	}

	err = tx.NewSelect().
		Model(&acc).
		Where("user_id = ?", userID).
		For("UPDATE").
		Scan(ctx)
	if err != nil {
		tx.Rollback() //nolint:errcheck
		return acc, err
	}

	if acc.CurrentPointsTotal < sum {
		tx.Rollback() //nolint:errcheck

		return acc, &store.NotEnoughBalanceError{
			Err:               errors.New("not enough balance"),
			Balance:           acc.CurrentPointsTotal,
			WithdrawRequested: sum,
		}
	}

	acc.CurrentPointsTotal = acc.CurrentPointsTotal - sum
	acc.WithdrawnTotal = acc.WithdrawnTotal + sum
	acc.UpdatedAt = time.Now()

	_, err = tx.NewUpdate().
		Model(&acc).
		WherePK().
		Column("current_points_total", "withdrawn_total", "updated_at").
		Returning("*").
		Exec(ctx)
	if err != nil {
		tx.Rollback() //nolint:errcheck
		return acc, err
	}

	transaction := models.Transaction{
		ID:          uuid.NewString(),
		AccountID:   acc.ID,
		Amount:      sum,
		OrderNumber: orderNumber,
		Direction:   models.TxDirectionWithdrawal,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = tx.NewInsert().
		Model(&transaction).
		Exec(ctx)
	if err != nil {
		tx.Rollback() //nolint:errcheck
		return acc, err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback() //nolint:errcheck
		return acc, err
	}

	return acc, nil
}

// GetWithdrawals fetches all transactions of specified direction
func (s *Store) GetWithdrawals(ctx context.Context, accountID string, direction models.TxDirection, limit int) ([]models.Transaction, error) {
	var result []models.Transaction

	err := s.conn.NewSelect().
		Model(&result).
		Where("account_id = ?", accountID).
		Where("direction = ?", direction).
		Order("created_at DESC").
		Limit(limit).
		Scan(ctx)

	return result, err
}

// AddBalance adds accrued points for order to account balance.
func (s *Store) AddBalance(ctx context.Context, order models.Order) (models.Account, error) {
	var acc models.Account

	tx, err := s.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return acc, err
	}

	err = tx.NewSelect().
		Model(&acc).
		Where("id = ?", order.AccountID).
		For("UPDATE").
		Scan(ctx)
	if err != nil {
		tx.Rollback() //nolint:errcheck
		return acc, err
	}

	acc.CurrentPointsTotal = acc.CurrentPointsTotal + order.Accrual
	acc.UpdatedAt = time.Now()

	_, err = tx.NewUpdate().
		Model(&acc).
		WherePK().
		Column("current_points_total", "updated_at").
		Returning("*").
		Exec(ctx)
	if err != nil {
		tx.Rollback() //nolint:errcheck
		return acc, err
	}

	transaction := models.Transaction{
		ID:          uuid.NewString(),
		AccountID:   acc.ID,
		Amount:      order.Accrual,
		OrderNumber: order.Number,
		Direction:   models.TxDirectionAccrual,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = tx.NewInsert().
		Model(&transaction).
		Exec(ctx)
	if err != nil {
		tx.Rollback() //nolint:errcheck
		return acc, err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback() //nolint:errcheck
		return acc, err
	}

	return acc, nil
}

func (s *Store) bootstrap(ctx context.Context) error {
	return Migrate(ctx, s.conn)
}

// Conn is used for migrations.
func (s *Store) Conn() *bun.DB {
	return s.conn
}
