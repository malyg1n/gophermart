package pgsql

import (
	"context"
	"gophermart/model"
	dbModel "gophermart/storage/pgsql/model"
)

// CreateOrder makes new order.
func (s Storage) CreateOrder(ctx context.Context, number string, userID uint64) error {
	_, err := s.db.ExecContext(
		ctx,
		"insert into orders (id, user_id) values ($1, $2);",
		number,
		userID,
	)

	return err
}

// GetOrderByNumber returns order by number.
func (s Storage) GetOrderByNumber(ctx context.Context, number string) (*model.Order, error) {
	var order dbModel.Order
	query := "select * from orders where id = $1"
	if err := s.db.GetContext(ctx, &order, query, number); err != nil {
		return nil, err
	}

	baseOrder := order.ToCanonical()

	return &baseOrder, nil
}

// GetOrdersByUser returns orders by user.
func (s Storage) GetOrdersByUser(ctx context.Context, userID uint64) ([]model.Order, error) {
	var orders []model.Order
	var dbOrders []dbModel.Order

	err := s.db.SelectContext(ctx, &dbOrders, "select * from orders where user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	for _, o := range dbOrders {
		orders = append(orders, o.ToCanonical())
	}

	return orders, nil
}

// UpdateOrder updates order.
func (s *Storage) UpdateOrder(ctx context.Context, number, status string, accrual int) error {
	_, err := s.db.ExecContext(
		ctx,
		"update orders set status = $1, accrual = $2 where id = $3;",
		status,
		accrual,
		number,
	)

	return err
}

// GetProcessOrders returns orders, which has status NEW or PROCESSING
func (s Storage) GetProcessOrders(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	var dbOrders []dbModel.Order

	err := s.db.SelectContext(ctx, &dbOrders, "select * from orders where status in ('NEW', 'PROCESSING')")
	if err != nil {
		return nil, err
	}

	for _, o := range dbOrders {
		orders = append(orders, o.ToCanonical())
	}

	return orders, nil
}
