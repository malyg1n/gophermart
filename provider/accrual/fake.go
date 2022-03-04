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
	var status string
	var accrual float64
	switch orderID {
	case "12345678903001":
		status = "PROCESSED"
		accrual = 500
	default:
		status = "PROCESSING"
		accrual = 0
	}

	return &model.Order{
		Number:  orderID,
		Status:  status,
		Accrual: accrual,
	}, nil
}
