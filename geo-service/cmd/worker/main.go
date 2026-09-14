package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

// Worker consumes GPS positions from RabbitMQ and persists them to PostgreSQL.
func main() {
	ctx := context.Background()
	app := fx.New(Module())

	if err := app.Start(ctx); err != nil {
		panic(fmt.Errorf("failed to start worker: %w", err))
	}

	sig := <-app.Wait()
	fmt.Printf("Worker stopped with code: %v\n", sig.ExitCode)

	if err := app.Stop(ctx); err != nil {
		fmt.Printf("error stopping worker: %v\n", err)
	}
}
