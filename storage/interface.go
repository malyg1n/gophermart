package storage

import (
	"context"
	"gophermart/model"
)

// UserStorer interface for operations with users.
type UserStorer interface {
	CreateUser(ctx context.Context, login, password string) error
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	GetUserByID(ctx context.Context, userID uint64) (*model.User, error)
}

// OrderStorer interface for operations with orders.
type OrderStorer interface {
	CreateOrder(ctx context.Context, number string, userID uint64) error
	GetOrderByNumber(ctx context.Context, number string) (*model.Order, error)
	GetOrdersByUser(ctx context.Context, userID uint64) ([]model.Order, error)
	UpdateOrder(ctx context.Context, number, status string, accrual int) error
}

// TransactionStorer interface for operations with transactions.
type TransactionStorer interface {
	GetOutcomeTransactionsByUser(ctx context.Context, userID uint64) ([]model.Transaction, error)
	SaveTransaction(ctx context.Context, userID uint64, orderID string, amount int) error
}
