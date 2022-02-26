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
	userStorage        storage.IUserStorage
	transactionStorage storage.ITransactionStorage
}

// NewUserService creates service instance.
func NewUserService(us storage.IUserStorage, ts storage.ITransactionStorage) UserService {
	return UserService{
		userStorage:        us,
		transactionStorage: ts,
	}
}

// Create user.
func (s UserService) Create(ctx context.Context, login, password string) error {
	password, err := s.cryptPassword(password)
	if err != nil {
		logger.GetLogger().Error(err)
		return err
	}
	_, err = s.userStorage.GetUserByLogin(ctx, login)
	if err == nil {
		return errs.ErrLoginExists
	}

	return s.userStorage.CreateUser(ctx, login, password)
}

// Auth user.
func (s UserService) Auth(ctx context.Context, login, password string) (string, error) {
	user, err := s.userStorage.GetUserByLogin(ctx, login)
	if err != nil {
		logger.GetLogger().Info(err)
		return "", errs.ErrAuthFailed
	}
	if err = s.comparePassword(password, user.CryptPassword); err != nil {
		logger.GetLogger().Info(err)
		return "", errs.ErrAuthFailed
	}

	return token.CreateTokenByUserID(user.ID)
}

// ShowBalance shows user balance.
func (s UserService) ShowBalance(ctx context.Context, userID int) (*response.Balance, error) {
	user, err := s.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return response.BalanceFromUser(user), nil
}

// GetTransactions by user
func (s UserService) GetTransactions(ctx context.Context, userID int) ([]*model.Transaction, error) {
	return s.transactionStorage.GetOutcomeTransactionsByUser(ctx, userID)
}

// TopUp user balance.
func (s UserService) TopUp(ctx context.Context, userID int, orderID string, amount float64) error {
	return s.transactionStorage.SaveTransaction(ctx, userID, orderID, amount)
}

// Withdraw money from user.
func (s UserService) Withdraw(ctx context.Context, userID int, orderID string, sum float64) error {
	user, err := s.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.Balance < sum {
		return errs.ErrBalanceTooSmall
	}

	sum = sum * -1

	return s.transactionStorage.SaveTransaction(ctx, userID, orderID, sum)
}

func (s UserService) cryptPassword(password string) (string, error) {
	crypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(crypt), nil
}

func (s UserService) comparePassword(password, cryptPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(cryptPassword), []byte(password))
}
