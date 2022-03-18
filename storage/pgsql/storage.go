package pgsql

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Storage base struct.
type Storage struct {
	db *sqlx.DB
}

// NewStorage returns new storage instance.
func NewStorage(dbURI string) (*Storage, error) {
	db, err := sqlx.Open("postgres", dbURI)
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
