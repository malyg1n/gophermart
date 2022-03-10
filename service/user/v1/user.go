package v1

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"gophermart/api/rest/response"
	"gophermart/model"
	"gophermart/pkg/errs"
	"gophermart/pkg/logger"
	"gophermart/pkg/token"
	"gophermart/storage"
)

// UserService for working with users.
type UserService struct {
	userStorage        storage.UserStorer
	transactionStorage storage.TransactionStorer
	logger             logger.Logger
}

type UserOption func(service *UserService)

// NewUserService constructor.
func NewUserService(opts ...UserOption) *UserService {
	service := &UserService{}

	for _, opt := range opts {
		opt(service)
	}

	return service
}

// WithUserStorageUserOption option.
func WithUserStorageUserOption(st storage.UserStorer) UserOption {
	return func(service *UserService) {
		service.userStorage = st
	}
}

// WithTransactionStorageUserOption option.
func WithTransactionStorageUserOption(st storage.TransactionStorer) UserOption {
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

// Create user.
func (s UserService) Create(ctx context.Context, login, password string) error {
	password, err := s.cryptPassword(password)
	if err != nil {
		s.logger.Errorf("%v", err)
		return err
	}
	user, _ := s.userStorage.GetUserByLogin(ctx, login)
	if user != nil {
		return errs.ErrLoginExists
	}

	return s.userStorage.CreateUser(ctx, login, password)
}

// Auth user.
func (s UserService) Auth(ctx context.Context, login, password string) (string, error) {
	user, err := s.userStorage.GetUserByLogin(ctx, login)
	if err != nil {
		s.logger.Infof("%v", err)
		return "", errs.ErrAuthFailed
	}
	if err = s.comparePassword(password, user.CryptPassword); err != nil {
		s.logger.Infof("%v", err)
		return "", errs.ErrAuthFailed
	}

	return token.CreateTokenByUserID(user.ID)
}

// ShowBalance shows user balance.
func (s UserService) ShowBalance(ctx context.Context, userID uint64) (*response.Balance, error) {
	user, err := s.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Errorf("%v", err)
		return nil, err
	}

	balance := response.BalanceFromUser(*user)

	return &balance, nil
}

// GetTransactions by user
func (s UserService) GetTransactions(ctx context.Context, userID uint64) ([]model.Transaction, error) {
	transactions, err := s.transactionStorage.GetOutcomeTransactionsByUser(ctx, userID)
	if err != nil {
		s.logger.Infof("%v", err)
	}
	return transactions, err
}

// Withdraw money from user.
func (s UserService) Withdraw(ctx context.Context, userID uint64, orderID string, sum float64) error {
	user, err := s.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Errorf("%v", err)
		return err
	}
	intSum := int(sum * 100)

	if user.Balance < intSum {
		return errs.ErrBalanceTooSmall
	}

	intSum = intSum * -1

	err = s.transactionStorage.SaveTransaction(ctx, userID, orderID, intSum)
	if err != nil {
		s.logger.Errorf("%v", err)
		return err
	}

	return nil
}

func (s UserService) cryptPassword(password string) (string, error) {
	crypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("%v", err)
		return "", err
	}

	return string(crypt), nil
}

func (s UserService) comparePassword(password, cryptPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(cryptPassword), []byte(password))
	if err != nil {
		s.logger.Errorf("%v", err)
	}

	return err
}
