package pgsql

import (
	"context"
	"fmt"
	"gophermart/model"
	dbModel "gophermart/storage/pgsql/model"
)

// CreateUser method inserts user in db.
func (s Storage) CreateUser(ctx context.Context, login, password string) error {
	_, err := s.db.ExecContext(
		ctx,
		"insert into users (login, password) values ($1, $2);",
		login,
		password,
	)

	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

// GetUserByLogin returns user by login
func (s Storage) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	var user dbModel.User
	query := "select * from users where login = $1"
	if err := s.db.GetContext(ctx, &user, query, login); err != nil {
		return nil, err
	}

	baseUser := user.ToCanonical()

	return &baseUser, nil
}

// GetUserByID returns user by login
func (s Storage) GetUserByID(ctx context.Context, id uint64) (*model.User, error) {
	var user dbModel.User
	query := "select * from users where id = $1"
	if err := s.db.GetContext(ctx, &user, query, id); err != nil {
		return nil, err
	}

	baseUser := user.ToCanonical()

	return &baseUser, nil
}
