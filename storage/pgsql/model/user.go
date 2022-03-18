package model

import "gophermart/model"

// User model.
type User struct {
	ID       uint64 `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
	Balance  int    `db:"balance"`
	Outcome  int    `db:"outcome"`
}

// ToCanonical converts db model to base model.
func (u User) ToCanonical() model.User {
	return model.User{
		ID:            u.ID,
		Login:         u.Login,
		CryptPassword: u.Password,
		Balance:       u.Balance,
		Outcome:       u.Outcome,
	}
}
