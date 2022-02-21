package response

import "gophermart/model"

// Order response.
type Order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAT string  `json:"uploaded_at"`
}

// OrderFromCanonical converts base model to response model.
func OrderFromCanonical(o *model.Order) Order {
	return Order{
		Number:     o.Number,
		Status:     o.Status,
		Accrual:    o.Accrual,
		UploadedAT: o.UploadedAT,
	}
}
