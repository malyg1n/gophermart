package model

import "gophermart/model"

// Order db model.
type Order struct {
	ID         string `db:"id"`
	UserID     uint64 `db:"user_id"`
	Status     string `db:"status"`
	Accrual    int    `db:"accrual"`
	UploadedAT string `db:"uploaded_at"`
}

// ToCanonical converts db model to base model.
func (o Order) ToCanonical() model.Order {
	return model.Order{
		Number:     o.ID,
		UserID:     o.UserID,
		Status:     o.Status,
		Accrual:    o.Accrual,
		UploadedAT: o.UploadedAT,
	}
}
