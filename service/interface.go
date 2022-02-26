package service

import (
	"context"
	"gophermart/api/rest/response"
	"gophermart/model"
)

// IUserService interface for operations with users.
type IUserService interface {
	Create(ctx context.Context, login, password string) error
	Auth(ctx context.Context, login, password string) (string, error)
	ShowBalance(ctx context.Context, userID int) (*response.Balance, error)
	GetTransactions(ctx context.Context, userID int) ([]*model.Transaction, error)
	Withdraw(ctx context.Context, userID int, orderID string, sum float64) error
	TopUp(ctx context.Context, userID int, orderID string, amount float64) error
}

// IOrderService interface for operations with orders.
type IOrderService interface {
	CreateOrder(ctx context.Context, number string, userID int) error
	GetOrderByNumber(ctx context.Context, number string) (*model.Order, error)
	GetOrdersByUser(ctx context.Context, userID int) ([]*model.Order, error)
	UpdateOrder(ctx context.Context, number, status string, accrual float64) error
}
