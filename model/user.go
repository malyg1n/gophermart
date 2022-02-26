package model

// User model.
type User struct {
	ID            int
	Login         string
	Password      string
	CryptPassword string
	Balance       float64
	Outcome       float64
}
