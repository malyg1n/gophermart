package v1

import (
	"gophermart/model"
	"gophermart/storage"
)

// Service for working with users.
type Service struct {
	userStorage storage.IUserStorage
}

// NewUserService creates service instance.
func NewUserService(userStorage storage.IUserStorage) Service {
	return Service{
		userStorage: userStorage,
	}
}

// Create user.
func (s Service) Create(login, password string) (*model.User, error) {
	return nil, nil
}

// Auth user.
func (s Service) Auth(login, password string) error {
	return nil
}
