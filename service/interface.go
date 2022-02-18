package service

import "gophermart/model"

// IUserService interface for operations with users.
type IUserService interface {
	Create(login, password string) (*model.User, error)
	Auth(login, password string) error
}
