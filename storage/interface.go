package storage

import (
	"context"
	"gophermart/model"
)

// IUserStorage interface for operations with users.
type IUserStorage interface {
	Create(ctx context.Context, login, password string) error
	GetByLogin(ctx context.Context, login string) (*model.User, error)
}
