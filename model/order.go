package model

// Order base model.
type Order struct {
	Number     string
	UserID     int
	Status     string
	Accrual    float64
	UploadedAT string
}
