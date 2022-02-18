package storage

import "gophermart/model"

// IUserStorage interface for operations with users.
type IUserStorage interface {
	Create(login, password string) error
	GetByLogin(login string) (*model.User, error)
}
