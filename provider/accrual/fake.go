package accrual

import (
	"gophermart/model"
)

// FakeHTTProvider struct.
type FakeHTTProvider struct {
}

// NewFakeHTTProvider struct.
func NewFakeHTTProvider() FakeHTTProvider {
	return FakeHTTProvider{}
}

// CheckOrder in accrual system.
func (p FakeHTTProvider) CheckOrder(orderID string) (*model.Order, error) {
	order := &model.Order{
		Number:  orderID,
		Status:  "PROCESSED",
		Accrual: 500,
	}

	return order, nil
}
