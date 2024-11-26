package database

import (
	"context"
	"embed"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

type Migrations struct {
	migrator *migrate.Migrator
}

//go:embed migrations/*.sql
var sqlMigrations embed.FS

// NewMigrations creates new Migrations.
func NewMigrations(conn *bun.DB) (*Migrations, error) {
	migrations := migrate.NewMigrations(migrate.WithMigrationsDirectory("internal/app/store/database/migrations"))
	if err := migrations.Discover(sqlMigrations); err != nil {
		return nil, err
	}

	return &Migrations{
		migrator: migrate.NewMigrator(conn, migrations),
	}, nil
}

// Migrate executes new migrations.
func (m *Migrations) Migrate(ctx context.Context) error {
	err := m.migrator.Init(ctx)
	if err != nil {
		return err
	}

	group, err := m.migrator.Migrate(ctx)
	if err != nil {
		return err
	}

	if group.ID == 0 {
		fmt.Printf("there are no new migrations to run\n")
		return nil
	}

	fmt.Printf("migrated to %s\n", group)
	return nil
}

// Rollback rollbacks last migration group
func (m *Migrations) Rollback(ctx context.Context) error {
	group, err := m.migrator.Rollback(ctx)
	if err != nil {
		return err
	}

	if group.ID == 0 {
		fmt.Printf("there are no groups to roll back\n")
		return nil
	}

	fmt.Printf("rolled back %s\n", group)
	return nil
}

// CreateMigration creates new migration.
func (m *Migrations) CreateMigration(ctx context.Context, name string) error {
	files, err := m.migrator.CreateSQLMigrations(ctx, name)
	if err != nil {
		return err
	}

	for _, mf := range files {
		fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
	}

	return nil
}
