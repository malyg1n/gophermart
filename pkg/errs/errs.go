package errs

import "errors"

var (
	ErrLoginExists          = errors.New("login already used")
	ErrAuthFailed           = errors.New("authentication failed")
	ErrOrderCreatedByMyself = errors.New("order already uploaded")
	ErrOrderExists          = errors.New("order already uploaded")
	ErrOrderNumber          = errors.New("invalid order number")
)
