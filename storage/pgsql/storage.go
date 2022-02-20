package pgsql

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gophermart/pkg/config"
)

// Storage base struct.
type Storage struct {
	db *sqlx.DB
}

// NewStorage returns new storage instance.
func NewStorage(cfg *config.AppConfig) (*Storage, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	db, err := sqlx.Open("postgres", cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	st := &Storage{
		db: db,
	}

	if err = st.Migrate(); err != nil {
		return nil, err
	}

	return st, nil
}
