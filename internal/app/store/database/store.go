package database

import (
	"context"

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

func (s *Store) bootstrap(ctx context.Context) error {
	return Migrate(ctx, s.conn)
}

// Conn is used for migrations.
func (s *Store) Conn() *bun.DB {
	return s.conn
}
