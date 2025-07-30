package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"binomeway.com/entry/application"
)

func main() {
	app := application.New(application.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		log.Fatal("failed to start app:", err)
	}
}
