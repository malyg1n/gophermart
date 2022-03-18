package model

// User model.
type User struct {
	ID            uint64
	Login         string
	Password      string
	CryptPassword string
	Balance       int
	Outcome       int
}
