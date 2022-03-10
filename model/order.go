package model

// Order base model.
type Order struct {
	Number     string
	UserID     uint64
	Status     string
	Accrual    int
	UploadedAT string
}
