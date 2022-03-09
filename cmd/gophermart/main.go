package main

import (
	"context"
	"gophermart/api/rest"
	"gophermart/api/rest/handler"
	"gophermart/pkg/config"
	"gophermart/pkg/logger"
	"gophermart/provider/accrual"
	v1 "gophermart/service/v1"
	"gophermart/storage/pgsql"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, ctxCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer ctxCancel()

	lgr := logger.NewDefaultLogger()
	cfg, err := config.NewDefaultConfig()
	if err != nil {
		lgr.Fatalf("config error %v", err)
		os.Exit(1)
	}

	accrualProvider := accrual.NewAccrualHTTPProvider(cfg.AccrualAddress)
	stg, err := pgsql.NewStorage(cfg.DatabaseURI)

	if err != nil {
		lgr.Fatalf("storage error %v", err)
		os.Exit(1)
	}

	userService := v1.NewUserService(
		v1.WithUserStorageUserOption(stg),
		v1.WithTransactionStorageUserOption(stg),
		v1.WithLoggerUserOption(lgr),
	)

	orderService := v1.NewOrderService(
		v1.WithOrderStorageOrderOption(stg),
		v1.WithTransactionStorageOrderOption(stg),
		v1.WithLoggerOrderOption(lgr),
		v1.WithProviderOrderOption(accrualProvider),
	)

	hr := handler.NewHandler(
		handler.WithUserService(userService),
		handler.WithOrderService(orderService),
		handler.WithLogger(lgr),
	)

	server := rest.NewAPIServer(hr, cfg.RunAddress)
	server.Run(ctx)

	<-ctx.Done()
	lgr.Infof("shutting down server")
}
