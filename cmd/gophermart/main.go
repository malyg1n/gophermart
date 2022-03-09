package main

import (
	"context"
	"gophermart/api/rest"
	"gophermart/api/rest/handler"
	"gophermart/pkg/config"
	"gophermart/pkg/logger"
	"gophermart/provider/accrual"
	orderService "gophermart/service/order/v1"
	userService "gophermart/service/user/v1"
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

	us := userService.NewUserService(
		userService.WithUserStorageUserOption(stg),
		userService.WithTransactionStorageUserOption(stg),
		userService.WithLoggerUserOption(lgr),
	)

	ose := orderService.NewOrderService(
		orderService.WithOrderStorageOrderOption(stg),
		orderService.WithTransactionStorageOrderOption(stg),
		orderService.WithLoggerOrderOption(lgr),
		orderService.WithProviderOrderOption(accrualProvider),
	)

	hr := handler.NewHandler(
		handler.WithUserService(us),
		handler.WithOrderService(ose),
		handler.WithLogger(lgr),
	)

	server := rest.NewAPIServer(hr, cfg.RunAddress)
	server.Run(ctx)

	<-ctx.Done()
	lgr.Infof("shutting down server")
}
