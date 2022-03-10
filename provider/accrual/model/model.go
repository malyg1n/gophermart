package model

import "gophermart/model"

// Order model.
type Order struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

// ToCanonical converts order to canonical model.
func (o Order) ToCanonical() *model.Order {
	return &model.Order{
		Number:  o.Order,
		Status:  o.Status,
		Accrual: o.Accrual,
	}
}
