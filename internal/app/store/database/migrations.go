package database

import (
	"context"
	"embed"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

var migrations = migrate.NewMigrations(migrate.WithMigrationsDirectory("internal/app/store/database/migrations"))

//go:embed migrations/*.sql
var sqlMigrations embed.FS

func init() {
	if err := migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
}

func Migrate(ctx context.Context, conn *bun.DB) error {
	migrator := migrate.NewMigrator(conn, migrations)

	err := migrator.Init(ctx)
	if err != nil {
		return err
	}

	group, err := migrator.Migrate(ctx)
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

func Rollback(ctx context.Context, conn *bun.DB) error {
	migrator := migrate.NewMigrator(conn, migrations)

	group, err := migrator.Rollback(ctx)
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

func CreateMigration(ctx context.Context, conn *bun.DB, name string) error {
	migrator := migrate.NewMigrator(conn, migrations)

	files, err := migrator.CreateSQLMigrations(ctx, name)
	if err != nil {
		return err
	}

	for _, mf := range files {
		fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
	}

	return nil
}
