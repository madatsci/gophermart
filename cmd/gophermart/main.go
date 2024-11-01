package main

import (
	"context"

	"github.com/madatsci/gophermart/internal/app"
	"github.com/madatsci/gophermart/internal/app/flags"
)

func main() {
	if err := flags.Parse(); err != nil {
		panic(err)
	}

	app, err := app.New(context.Background(), app.Options{
		RunAddress:           flags.RunAddress,
		AccrualSystemAddress: flags.AccrualSystemAddress,
		DatabaseURI:          flags.DatabaseURI,
	})
	if err != nil {
		panic(err)
	}

	if err = app.Start(); err != nil {
		panic(err)
	}
}
