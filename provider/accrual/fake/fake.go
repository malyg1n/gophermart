package fake

import (
	"gophermart/model"
)

// HTTPProvider struct.
type HTTPProvider struct {
}

// NewFakeHTTProvider struct.
func NewFakeHTTProvider() HTTPProvider {
	return HTTPProvider{}
}

// CheckOrder in accrual system.
func (p HTTPProvider) CheckOrder(orderID string) (*model.Order, error) {
	var status string
	var accrual float64
	switch orderID {
	case "12345678903001":
		status = "PROCESSED"
		accrual = 500.28
	default:
		status = "PROCESSING"
		accrual = 0
	}

	return &model.Order{
		Number:  orderID,
		Status:  status,
		Accrual: int(accrual * 100),
	}, nil
}
