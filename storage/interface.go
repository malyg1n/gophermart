package storage

import (
	"context"
	"gophermart/model"
)

// IUserStorage interface for operations with users.
type IUserStorage interface {
	CreateUser(ctx context.Context, login, password string) error
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
}

// IOrderStorage interface for operations with orders.
type IOrderStorage interface {
	CreateOrder(ctx context.Context, number string, userID int) error
	GetOrderByNumber(ctx context.Context, number string) (*model.Order, error)
	GetOrdersByUser(ctx context.Context, userID int) ([]*model.Order, error)
	UpdateOrder(ctx context.Context, order model.Order) error
}
