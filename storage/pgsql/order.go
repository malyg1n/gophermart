package pgsql

import (
	"context"
	"gophermart/model"
	dbModel "gophermart/storage/pgsql/model"
)

func (s *Storage) CreateOrder(ctx context.Context, number string, userID int) error {
	_, err := s.db.ExecContext(
		ctx,
		"insert into orders (id, user_id) values ($1, $2);",
		number,
		userID,
	)

	return err
}

func (s *Storage) GetOrderByNumber(ctx context.Context, number string) (*model.Order, error) {
	var order dbModel.Order
	query := "select * from orders where id = $1"
	if err := s.db.GetContext(ctx, &order, query, number); err != nil {
		return nil, err
	}

	return order.ToCanonical(), nil
}

func (s *Storage) GetOrdersByUser(ctx context.Context, userID int) ([]*model.Order, error) {
	var orders []*model.Order
	var dbOrders []*dbModel.Order

	err := s.db.SelectContext(ctx, &dbOrders, "select * from orders where user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	for _, o := range dbOrders {
		orders = append(orders, o.ToCanonical())
	}

	return orders, nil
}

func (s *Storage) UpdateOrder(ctx context.Context, order model.Order) error {
	_, err := s.db.ExecContext(
		ctx,
		"update orders set status = $1, accrual = $2 where id = $3;",
		order.Status,
		order.Accrual,
		order.Number,
	)

	return err
}
