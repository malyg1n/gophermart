package v1

import (
	"context"
	"errors"
	"gophermart/model"
	"gophermart/pkg/config"
	"gophermart/pkg/errs"
	"gophermart/pkg/logger"
	"gophermart/pkg/validation"
	"gophermart/provider/accrual"
	"gophermart/storage"
	"time"
)

// OrderService struct.
type OrderService struct {
	orderStorage       storage.IOrderStorage
	transactionStorage storage.ITransactionStorage
}

// NewOrderService constructor.
func NewOrderService(os storage.IOrderStorage, ts storage.ITransactionStorage) OrderService {
	return OrderService{
		orderStorage:       os,
		transactionStorage: ts,
	}
}

// CreateOrder makes new order.
func (s OrderService) CreateOrder(ctx context.Context, number string, userID int) error {
	if number == "" || validation.IsLunh(number) == false {
		return errs.ErrOrderNumber
	}

	order, err := s.orderStorage.GetOrderByNumber(ctx, number)
	if err == nil {
		if order.UserID == userID {
			return errs.ErrOrderCreatedByMyself
		}
		return errs.ErrOrderExists
	}

	go func() {
		err = s.processOrder(ctx, number)
		if err != nil {
			logger.GetLogger().Error(err)
		}
	}()

	return s.orderStorage.CreateOrder(ctx, number, userID)
}

// GetOrderByNumber returns order by number.
func (s OrderService) GetOrderByNumber(ctx context.Context, number string) (*model.Order, error) {
	return s.orderStorage.GetOrderByNumber(ctx, number)
}

// GetOrdersByUser returns orders by user.
func (s OrderService) GetOrdersByUser(ctx context.Context, userID int) ([]*model.Order, error) {
	return s.orderStorage.GetOrdersByUser(ctx, userID)
}

// updateOrder updates order.
func (s OrderService) updateOrder(ctx context.Context, number, status string, accrual float64) error {
	return s.orderStorage.UpdateOrder(ctx, number, status, accrual)
}

// processOrder check order in accrual system.
func (s OrderService) processOrder(ctx context.Context, orderID string) error {
	var i int
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}
	provider := accrual.NewAccrualHttpProvider(cfg.AccrualAddress)
	status := "NEW"
	// 100 tries to check order
	for i < 100 {
		order, err := provider.CheckOrder(orderID)
		if err != nil && errors.Is(errs.ErrToManyRequests, err) {
			time.Sleep(time.Second * 60)
			continue
		}
		if status != order.Status {
			status = order.Status
			err = s.updateOrder(ctx, order.Number, order.Status, order.Accrual)
			if err != nil {
				return err
			}
		}

		if status == "PROCESSED" || status == "INVALID" {
			return nil
		}

		time.Sleep(time.Second)
		i++
	}

	return nil
}
