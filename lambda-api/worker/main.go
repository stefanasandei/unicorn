package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"binomeway.com/worker/application"
)

func main() {
	app := application.New(application.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		log.Panicf("Failed to start the app: %s", err)
	}
}
