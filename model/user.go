package model

// User model.
type User struct {
	ID            uint64
	Login         string
	Password      string
	CryptPassword string
	Balance       float64
	Outcome       float64
}
