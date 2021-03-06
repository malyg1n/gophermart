package v1

import (
	"context"
	"gophermart/model"
	"gophermart/pkg/errs"
	"gophermart/pkg/logger"
	"gophermart/pkg/validation"
	"gophermart/provider"
	"gophermart/storage"
	"time"
)

// OrderService struct.
type OrderService struct {
	orderStorage       storage.OrderStorer
	transactionStorage storage.TransactionStorer
	logger             logger.Logger
	provider           provider.AccrualProvider
}

type OrderOption func(service *OrderService)

// NewOrderService constructor.
func NewOrderService(opts ...OrderOption) *OrderService {
	service := &OrderService{}

	for _, opt := range opts {
		opt(service)
	}

	return service
}

// WithOrderStorageOrderOption option.
func WithOrderStorageOrderOption(st storage.OrderStorer) OrderOption {
	return func(service *OrderService) {
		service.orderStorage = st
	}
}

// WithTransactionStorageOrderOption option.
func WithTransactionStorageOrderOption(st storage.TransactionStorer) OrderOption {
	return func(service *OrderService) {
		service.transactionStorage = st
	}
}

// WithLoggerOrderOption option.
func WithLoggerOrderOption(l logger.Logger) OrderOption {
	return func(service *OrderService) {
		service.logger = l
	}
}

// WithProviderOrderOption option.
func WithProviderOrderOption(p provider.AccrualProvider) OrderOption {
	return func(service *OrderService) {
		service.provider = p
	}
}

// CreateOrder makes new order.
func (s OrderService) CreateOrder(ctx context.Context, number string, userID uint64) error {
	if !validation.IsLunh(number) {
		return errs.ErrOrderNumber
	}

	order, err := s.orderStorage.GetOrderByNumber(ctx, number)
	if err == nil {
		if order.UserID == userID {
			return errs.ErrOrderCreatedByMyself
		}
		return errs.ErrOrderExists
	}

	go s.processOrder(number, userID, "NEW")

	return s.orderStorage.CreateOrder(ctx, number, userID)
}

// GetOrderByNumber returns order by number.
func (s OrderService) GetOrderByNumber(ctx context.Context, number string) (*model.Order, error) {
	orders, err := s.orderStorage.GetOrderByNumber(ctx, number)
	if err != nil {
		s.logger.Errorf("%v", err)
	}

	return orders, err
}

// GetOrdersByUser returns orders by user.
func (s OrderService) GetOrdersByUser(ctx context.Context, userID uint64) ([]model.Order, error) {
	orders, err := s.orderStorage.GetOrdersByUser(ctx, userID)
	if err != nil {
		s.logger.Errorf("%v", err)
	}

	return orders, err
}

// updateOrder updates order.
func (s *OrderService) updateOrder(ctx context.Context, number, status string, accrual int) error {
	err := s.orderStorage.UpdateOrder(ctx, number, status, accrual)
	if err != nil {
		s.logger.Errorf("%v", err)
	}

	return err
}

// ProcessOrders check orders in accrual system.
func (s OrderService) ProcessOrders() {
	orders, err := s.orderStorage.GetProcessOrders(context.Background())
	if err != nil {
		s.logger.Errorf("%v", err)
		return
	}

	for _, order := range orders {
		go s.processOrder(order.Number, order.UserID, order.Status)
	}
}

// processOrder check order in accrual system.
func (s OrderService) processOrder(orderID string, userID uint64, status string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	order, err := s.provider.CheckOrder(orderID)
	if err != nil {
		s.logger.Errorf("%v", err)
		return
	}

	if status != order.Status {
		status = order.Status
		s.logger.Infof("update order %s, status=%s, accrual=%v", order.Number, order.Status, order.Accrual)
		err = s.updateOrder(ctx, order.Number, order.Status, order.Accrual)
		if err != nil {
			s.logger.Errorf("%v", err)
			return
		}
	}

	if status == "PROCESSED" || status == "INVALID" {
		if status == "PROCESSED" {
			err = s.transactionStorage.SaveTransaction(ctx, userID, orderID, order.Accrual)
			if err != nil {
				s.logger.Errorf("%v", err)
			}
		}
		return
	}
}
