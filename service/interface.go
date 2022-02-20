package service

import (
	"context"
)

// IUserService interface for operations with users.
type IUserService interface {
	Create(ctx context.Context, login, password string) error
	Auth(ctx context.Context, login, password string) (string, error)
}
