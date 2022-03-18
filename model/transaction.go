package model

// Transaction base model.
type Transaction struct {
	UserID    string
	OrderID   string
	Amount    int
	CreatedAt string
}
