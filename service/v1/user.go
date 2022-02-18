package v1

import (
	"gophermart/model"
	"gophermart/storage"
)

// UserService for working with users.
type UserService struct {
	userStorage storage.IUserStorage
}

// NewUserService creates service instance.
func NewUserService(userStorage storage.IUserStorage) UserService {
	return UserService{
		userStorage: userStorage,
	}
}

// Create user.
func (s UserService) Create(login, password string) (*model.User, error) {
	return nil, nil
}

// Auth user.
func (s UserService) Auth(login, password string) error {
	return nil
}
