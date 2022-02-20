package v1

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"gophermart/pkg/errs"
	"gophermart/pkg/logger"
	"gophermart/pkg/token"
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
