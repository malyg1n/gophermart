package model

import "gophermart/model"

// User model.
type User struct {
	ID       int     `db:"id"`
	Login    string  `db:"login"`
	Password string  `db:"password"`
	Balance  float64 `db:"balance"`
	Outcome  float64 `db:"outcome"`
}

// ToCanonical converts db model to base model.
func (u User) ToCanonical() *model.User {
	return &model.User{
		ID:            u.ID,
		Login:         u.Login,
		CryptPassword: u.Password,
	}
}
