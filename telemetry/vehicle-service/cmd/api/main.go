package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

// @title           Vehicle Service API
// @version         1.0
// @description     Vehicle management microservice for the telemetry platform.
// @host            localhost:8080
// @BasePath        /
// @schemes         http
func main() {
	ctx := context.Background()
	app := fx.New(Module())

	if err := app.Start(ctx); err != nil {
		panic(fmt.Errorf("failed to start application: %w", err))
	}

	sig := <-app.Wait()
	fmt.Printf("Application stopped with code: %v\n", sig.ExitCode)

	if err := app.Stop(ctx); err != nil {
		fmt.Printf("error stopping application: %v\n", err)
	}
}
