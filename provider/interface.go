package provider

import "gophermart/model"

// AccrualProvider provider interface.
type AccrualProvider interface {
	CheckOrder(orderID string) (*model.Order, error)
}
