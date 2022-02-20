package errs

import "errors"

var (
	ErrLoginExists = errors.New("login already used")
	ErrAuthFailed  = errors.New("authentication failed")
)
