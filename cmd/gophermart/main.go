package main

import (
	"context"
	"gophermart/api/rest"
	"gophermart/api/rest/handler"
	"gophermart/pkg/config"
	"gophermart/pkg/logger"
	accrualHTTPProvider "gophermart/provider/accrual/http"
	orderService "gophermart/service/order/v1"
	userService "gophermart/service/user/v1"
	"gophermart/storage/pgsql"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	accrualProvider := accrualHTTPProvider.NewAccrualHTTPProvider(cfg.AccrualAddress)
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

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			<-ticker.C
			ose.ProcessOrders()
		}
	}()

	<-ctx.Done()
	ticker.Stop()
	lgr.Infof("shutting down server")
}
