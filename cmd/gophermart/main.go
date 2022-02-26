package main

import (
	"context"
	"gophermart/api/rest"
	"gophermart/api/rest/handler"
	"gophermart/pkg/config"
	"gophermart/pkg/logger"
	v1 "gophermart/service/v1"
	"gophermart/storage/pgsql"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, ctxCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer ctxCancel()

	cfg, err := config.GetConfig()
	if err != nil {
		logger.GetLogger().Errorw("config error", "error", err.Error())
		os.Exit(1)
	}

	stg, err := pgsql.NewStorage(cfg)
	if err != nil {
		logger.GetLogger().Errorw("storage error", "error", err.Error())
		os.Exit(1)
	}

	userService := v1.NewUserService(stg, stg)
	orderService := v1.NewOrderService(stg, stg)

	hr := handler.NewHandler(
		handler.WithUserService(userService),
		handler.WithOrderService(orderService),
	)

	server := rest.NewAPIServer(hr, cfg.RunAddress)
	server.Run(ctx)

	<-ctx.Done()
	logger.GetLogger().Info("shutting down server")
}
