package model

import "gophermart/model"

// Transaction db model.
type Transaction struct {
	ID        int     `db:"id"`
	UserID    string  `db:"user_id"`
	OrderID   string  `db:"order_id"`
	Amount    float64 `db:"amount"`
	CreatedAt string  `db:"created_at"`
}

// ToCanonical converts db model to base model.
func (t Transaction) ToCanonical() *model.Transaction {
	return &model.Transaction{
		UserID:    t.UserID,
		OrderID:   t.OrderID,
		Amount:    t.Amount,
		CreatedAt: t.CreatedAt,
	}
}
