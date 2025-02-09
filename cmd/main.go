//go:generate oapi-codegen -generate types,models,gin -package api -o ../pkg/api/api.gen.go ../api/openapi.yaml
package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Rolan335/Finop/internal/app"
	"github.com/Rolan335/Finop/internal/config"
	"github.com/Rolan335/Finop/internal/controller"
	"github.com/Rolan335/Finop/internal/repository/postgres"
	"github.com/Rolan335/Finop/internal/service/finop"
)

func main() {
	envPath := ".env"
	cfg := config.MustNewConfig(envPath)
	if err := postgres.Migrate(&cfg.Migration); err != nil {
		panic("failed to migrate: " + err.Error())
	}
	storage := postgres.MustNewStorage(&cfg.DB)

	service := finop.NewFinop(storage)

	server := controller.New(cfg.RequestTimeout, service)

	app := app.NewService(cfg, server)

	//creating notify ctx for graceful shutdown and starting app
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	app.Start()
	<-ctx.Done()

	//stopping server and provided services. Provided servies should have method Stop()
	app.GracefulStop(storage)
}
