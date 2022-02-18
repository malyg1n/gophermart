package pgsql

import "gophermart/pkg/config"

// Storage base struct.
type Storage struct {
}

// NewStorage returns new storage instance.
func NewStorage(cfg *config.AppConfig) (*Storage, error) {
	return &Storage{}, nil
}
