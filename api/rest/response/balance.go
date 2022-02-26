package response

import "gophermart/model"

// Balance json model.
type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

// BalanceFromUser converts base user model to Balance.
func BalanceFromUser(u *model.User) *Balance {
	return &Balance{
		Current:   u.Balance,
		Withdrawn: u.Outcome,
	}
}
