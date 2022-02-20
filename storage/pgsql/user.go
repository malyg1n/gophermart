package pgsql

import (
	"context"
	"gophermart/model"
	dbModel "gophermart/storage/pgsql/model"
)

// Create method inserts user in db.
func (s *Storage) Create(ctx context.Context, login, password string) error {
	_, err := s.db.ExecContext(
		ctx,
		"insert into users (login, password) values ($1, $2);",
		login,
		password,
	)

	return err
}

// GetByLogin returns user by login
func (s *Storage) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	var user dbModel.User
	query := "select * from users where login = $1"
	if err := s.db.GetContext(ctx, &user, query, login); err != nil {
		return nil, err
	}

	return user.ToCanonical(), nil
}
