package errs

import "errors"

var (
	ErrLoginExists          = errors.New("login already used")
	ErrAuthFailed           = errors.New("authentication failed")
	ErrOrderCreatedByMyself = errors.New("order already uploaded")
	ErrOrderExists          = errors.New("order already uploaded")
	ErrOrderNumber          = errors.New("invalid order number")
	ErrBalanceTooSmall      = errors.New("balance too small")
	ErrToManyRequests       = errors.New("too many request")
	ErrAccrualResponse      = errors.New("accrual system is not available")
)
