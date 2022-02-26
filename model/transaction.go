package model

// Transaction base model.
type Transaction struct {
	UserID    string
	OrderID   string
	Amount    float64
	CreatedAt string
}
