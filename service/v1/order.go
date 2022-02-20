package v1

import (
	"context"
	"gophermart/model"
	"gophermart/pkg/errs"
	"gophermart/pkg/validation"
	"gophermart/storage"
)

type OrderService struct {
	storage storage.IOrderStorage
}

func NewOrderService(st storage.IOrderStorage) OrderService {
	return OrderService{
		storage: st,
	}
}

func (s OrderService) CreateOrder(ctx context.Context, number string, userID int) error {
	if !validation.IsLunh(number) {
		return errs.ErrOrderNumber
	}

	order, err := s.storage.GetOrderByNumber(ctx, number)
	if err == nil {
		if order.UserID == userID {
			return errs.ErrOrderCreatedByMyself
		}
		return errs.ErrOrderExists
	}

	return s.storage.CreateOrder(ctx, number, userID)
}

func (s OrderService) GetOrderByNumber(ctx context.Context, number string) (*model.Order, error) {
	return s.storage.GetOrderByNumber(ctx, number)
}

func (s OrderService) GetOrdersByUser(ctx context.Context, userID int) ([]*model.Order, error) {
	return s.storage.GetOrdersByUser(ctx, userID)
}

func (s OrderService) UpdateOrder(ctx context.Context, number, status string, accrual float64) error {
	return nil
}
