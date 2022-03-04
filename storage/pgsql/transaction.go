package pgsql

import (
	"context"
	"gophermart/model"
	dbModel "gophermart/storage/pgsql/model"
)

// GetOutcomeTransactionsByUser returns outcome transactions bu user.
func (s *Storage) GetOutcomeTransactionsByUser(ctx context.Context, userID int) ([]*model.Transaction, error) {
	dbTrans := make([]*dbModel.Transaction, 0)

	err := s.db.SelectContext(ctx, &dbTrans, "select * from transactions where user_id = $1 and amount < 0", userID)
	if err != nil {
		return nil, err
	}

	trans := make([]*model.Transaction, 0, len(dbTrans))
	for _, t := range dbTrans {
		trans = append(trans, t.ToCanonical())
	}

	return trans, nil
}

// SaveTransaction creates new transaction
func (s *Storage) SaveTransaction(ctx context.Context, userID int, orderID string, amount float64) error {
	_, err := s.db.ExecContext(
		ctx,
		"insert into transactions (user_id, order_id, amount) values ($1, $2, $3)",
		userID,
		orderID,
		amount,
	)

	if err != nil {
		return err
	}

	return s.updateUserBalance(ctx, userID, amount)
}

// updateUserBalance update user balance.
func (s *Storage) updateUserBalance(ctx context.Context, userID int, amount float64) error {
	_, err := s.db.ExecContext(
		ctx,
		"update users set balance = balance + $1 where id = $2",
		amount,
		userID,
	)
	if err != nil {
		return err
	}

	if amount < 0 {
		_, err := s.db.ExecContext(
			ctx,
			"update users set outcome = outcome - $1 where id = $2",
			amount,
			userID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
