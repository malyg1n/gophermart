package pgsql

import "gophermart/model"

// Create method inserts user in db.
func (s *Storage) Create(login, password string) error {
	return nil
}

// GetByLogin returns user by login
func (s *Storage) GetByLogin(login string) (*model.User, error) {
	return nil, nil
}
