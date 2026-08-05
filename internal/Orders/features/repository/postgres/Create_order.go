package orders_repository_posgres

import (
	"context"
	orders_domains "couriers/internal/Orders/core/domains"
	"fmt"
)

func (d *DatabasePostgres) CreateOrder(
	ctx context.Context,
	order orders_domains.Order,
) (orders_domains.Order, error) {
	ctxWithTime, cancel := context.WithTimeout(ctx, d.pool.GetTimeot())
	defer cancel()

	sqlRequest := `
	INSERT INTO app.orders (name,price,is_complete,user_id,courier_id)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id, version, name, price, is_complete, user_id, courier_id;
	`

	row := d.pool.QueryRow(
		ctxWithTime,
		sqlRequest,
		order.Name,
		order.Price,
		order.IsComplete,
		order.UserID,
		order.CourierID,
	)

	var orderModel OrderModel

	if err := row.Scan(
		&orderModel.ID,
		&orderModel.Version,
		&orderModel.Name,
		&orderModel.Price,
		&orderModel.IsComplete,
		&orderModel.UserID,
		&orderModel.CourierID,
	); err != nil {
		return orders_domains.Order{}, fmt.Errorf("fail to scan order: %w", err)
	}

	newOrder := orderDomainFromModel(orderModel)

	return newOrder, nil
}
