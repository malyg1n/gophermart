package accrual

import "gophermart/model"

// IAccrualProvider provider interface.
type IAccrualProvider interface {
	CheckOrder(orderID string) (*model.Order, error)
}
