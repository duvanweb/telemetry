package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

// @title           Alert Service API
// @version         1.0
// @description     Alert management microservice for the telemetry platform.
// @host            localhost:8082
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
