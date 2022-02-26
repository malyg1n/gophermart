package response

import "gophermart/model"

// Transaction response model.
type Transaction struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

// TransactionFromCanonical converts base model to response model.
func TransactionFromCanonical(t *model.Transaction) Transaction {
	return Transaction{
		Order:       t.OrderID,
		Sum:         t.Amount * -1.0,
		ProcessedAt: t.CreatedAt,
	}
}
