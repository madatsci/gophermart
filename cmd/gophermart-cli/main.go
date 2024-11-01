package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/madatsci/gophermart/internal/app"
	"github.com/madatsci/gophermart/internal/app/store/database"
	"github.com/uptrace/bun"
	"github.com/urfave/cli/v2"
)

var (
	Name    = "gophermart-cli"
	Version = "development"
)

func main() {
	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "d",
			Usage: "database URI",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "db",
			Usage: "manage database migrations",
			Subcommands: []*cli.Command{
				{
					Name:  "create_sql",
					Usage: "create SQL transaction",
					Action: func(cliCtx *cli.Context) error {
						return run(cliCtx, func(ctx context.Context, conn *bun.DB) error {
							name := strings.Join(cliCtx.Args().Slice(), "_")
							return database.CreateMigration(ctx, conn, name)
						})
					},
				},
				{
					Name:  "migrate",
					Usage: "migrate database",
					Action: func(cliCtx *cli.Context) error {
						return run(cliCtx, func(ctx context.Context, conn *bun.DB) error {
							return database.Migrate(ctx, conn)
						})
					},
				},
				{
					Name:  "rollback",
					Usage: "rollback the last migration group",
					Action: func(cliCtx *cli.Context) error {
						return run(cliCtx, func(ctx context.Context, conn *bun.DB) error {
							return database.Rollback(ctx, conn)
						})
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(cliCtx *cli.Context, fn func(ctx context.Context, conn *bun.DB) error) error {
	ctx := cliCtx.Context

	app, err := app.New(ctx, app.Options{
		DatabaseURI: cliCtx.String("d"),
	})
	if err != nil {
		return err
	}

	store, ok := app.Store().(*database.Store)
	if !ok {
		return errors.New("database store must be used for migrations")
	}

	return fn(ctx, store.Conn())
}
