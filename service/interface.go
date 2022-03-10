package service

import (
	"context"
	"gophermart/api/rest/response"
	"gophermart/model"
)

// UserProcessor interface for operations with users.
type UserProcessor interface {
	Create(ctx context.Context, login, password string) error
	Auth(ctx context.Context, login, password string) (string, error)
	ShowBalance(ctx context.Context, userID uint64) (*response.Balance, error)
	GetTransactions(ctx context.Context, userID uint64) ([]model.Transaction, error)
	Withdraw(ctx context.Context, userID uint64, orderID string, sum float64) error
}

// OrderProcessor interface for operations with orders.
type OrderProcessor interface {
	CreateOrder(ctx context.Context, number string, userID uint64) error
	GetOrderByNumber(ctx context.Context, number string) (*model.Order, error)
	GetOrdersByUser(ctx context.Context, userID uint64) ([]model.Order, error)
}
