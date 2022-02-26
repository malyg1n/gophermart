package request

// Withdraw request model
type Withdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
