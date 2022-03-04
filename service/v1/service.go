package v1

import (
	"gophermart/pkg/logger"
	"gophermart/provider/accrual"
	"gophermart/storage"
)

// UserService for working with users.
type UserService struct {
	userStorage        storage.IUserStorage
	transactionStorage storage.ITransactionStorage
	logger             logger.Logger
}

// OrderService struct.
type OrderService struct {
	orderStorage       storage.IOrderStorage
	transactionStorage storage.ITransactionStorage
	logger             logger.Logger
	provider           accrual.IAccrualProvider
}

type OrderOption func(service *OrderService)
type UserOption func(service *UserService)

// NewOrderService constructor.
func NewOrderService(opts ...OrderOption) *OrderService {
	service := &OrderService{}

	for _, opt := range opts {
		opt(service)
	}

	return service
}

// NewUserService constructor.
func NewUserService(opts ...UserOption) *UserService {
	service := &UserService{}

	for _, opt := range opts {
		opt(service)
	}

	return service
}

// WithOrderStorageOrderOption option.
func WithOrderStorageOrderOption(st storage.IOrderStorage) OrderOption {
	return func(service *OrderService) {
		service.orderStorage = st
	}
}

// WithTransactionStorageOrderOption option.
func WithTransactionStorageOrderOption(st storage.ITransactionStorage) OrderOption {
	return func(service *OrderService) {
		service.transactionStorage = st
	}
}

// WithLoggerOrderOption option.
func WithLoggerOrderOption(l logger.Logger) OrderOption {
	return func(service *OrderService) {
		service.logger = l
	}
}

// WithProviderOrderOption option.
func WithProviderOrderOption(p accrual.IAccrualProvider) OrderOption {
	return func(service *OrderService) {
		service.provider = p
	}
}

// WithUserStorageUserOption option.
func WithUserStorageUserOption(st storage.IUserStorage) UserOption {
	return func(service *UserService) {
		service.userStorage = st
	}
}

// WithTransactionStorageUserOption option.
func WithTransactionStorageUserOption(st storage.ITransactionStorage) UserOption {
	return func(service *UserService) {
		service.transactionStorage = st
	}
}

// WithLoggerUserOption option.
func WithLoggerUserOption(l logger.Logger) UserOption {
	return func(service *UserService) {
		service.logger = l
	}
}
