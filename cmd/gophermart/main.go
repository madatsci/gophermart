package main

import (
	"context"

	"github.com/madatsci/gophermart/internal/app"
)

func main() {
	if err := parseFlags(); err != nil {
		panic(err)
	}

	app, err := app.New(context.Background(), app.Options{
		RunAddress:           runAddress,
		AccrualSystemAddress: accrualSystemAddress,
		DatabaseURI:          databaseURI,
	})
	if err != nil {
		panic(err)
	}

	if err = app.Start(); err != nil {
		panic(err)
	}
}
