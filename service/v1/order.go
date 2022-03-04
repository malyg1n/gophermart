package v1

import (
	"context"
	"errors"
	"gophermart/model"
	"gophermart/pkg/errs"
	"gophermart/pkg/validation"
	"time"
)

// CreateOrder makes new order.
func (s *OrderService) CreateOrder(ctx context.Context, number string, userID int) error {
	if !validation.IsLunh(number) {
		return errs.ErrOrderNumber
	}

	order, err := s.orderStorage.GetOrderByNumber(ctx, number)
	if err == nil {
		s.logger.Errorf("%v", err)
		if order.UserID == userID {
			return errs.ErrOrderCreatedByMyself
		}
		return errs.ErrOrderExists
	}

	go s.processOrder(number, userID)

	return s.orderStorage.CreateOrder(ctx, number, userID)
}

// GetOrderByNumber returns order by number.
func (s *OrderService) GetOrderByNumber(ctx context.Context, number string) (*model.Order, error) {
	orders, err := s.orderStorage.GetOrderByNumber(ctx, number)
	if err != nil {
		s.logger.Errorf("%v", err)
	}

	return orders, err
}

// GetOrdersByUser returns orders by user.
func (s OrderService) GetOrdersByUser(ctx context.Context, userID int) ([]*model.Order, error) {
	orders, err := s.orderStorage.GetOrdersByUser(ctx, userID)
	if err != nil {
		s.logger.Errorf("%v", err)
	}

	return orders, err
}

// updateOrder updates order.
func (s *OrderService) updateOrder(ctx context.Context, number, status string, accrual float64) error {
	err := s.orderStorage.UpdateOrder(ctx, number, status, accrual)
	if err != nil {
		s.logger.Errorf("%v", err)
	}

	return err
}

// processOrder check order in accrual system.
func (s *OrderService) processOrder(orderID string, userID int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	var i int
	status := "NEW"
	// 100 tries to check order (@todo refactor this shit)
	for i < 100 {
		order, err := s.provider.CheckOrder(orderID)
		if err != nil {
			if errors.Is(errs.ErrToManyRequests, err) {
				time.Sleep(time.Second * 60)
				continue
			}
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

		time.Sleep(time.Second)
		i++
	}
}
