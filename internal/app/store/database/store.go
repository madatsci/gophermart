package database

import (
	"context"

	"github.com/madatsci/gophermart/internal/app/models"
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

func (s *Store) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	var result models.User

	err := s.conn.NewInsert().Model(&user).Returning("*").Scan(ctx, &result)

	return result, err
}

func (s *Store) bootstrap(ctx context.Context) error {
	return Migrate(ctx, s.conn)
}

// Conn is used for migrations.
func (s *Store) Conn() *bun.DB {
	return s.conn
}
