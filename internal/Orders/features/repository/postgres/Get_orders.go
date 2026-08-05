package orders_repository_posgres

import (
	"context"
	orders_domains "couriers/internal/Orders/core/domains"
	"fmt"
)

func (d *DatabasePostgres) GetOrders(
	ctx context.Context,
	personID int,
	limit *int,
	offset *int,
) ([]orders_domains.GetOrder, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, d.pool.GetTimeot())
	defer cancel()

	sqlRequest := `
	SELECT app.orders.id, version, name, price, is_complete, user_id, courier_id, app.couriers.login 
	FROM app.orders LEFT JOIN app.couriers ON courier_id=app.couriers.id LEFT JOIN app.users 
	ON user_id = app.users.id
	WHERE app.users.id = $1;
	`

	rows, err := d.pool.Query(ctxWithTimeout, sqlRequest, personID)
	if err != nil {
		return nil, fmt.Errorf("get orders from database error: %w", err)
	}
	defer rows.Close()

	var orderModels []GetOrderModel

	for rows.Next() {
		var orderModel GetOrderModel
		if err := rows.Scan(
			&orderModel.ID,
			&orderModel.Version,
			&orderModel.Name,
			&orderModel.Price,
			&orderModel.IsComplete,
			&orderModel.UserID,
			&orderModel.CourierID,
			&orderModel.CourierLogin,
		); err != nil {
			return nil, fmt.Errorf("get order error from database: %w", err)
		}

		orderModels = append(orderModels, orderModel)
	}

	orders := make([]orders_domains.GetOrder, len(orderModels))

	for i, model := range orderModels {
		orders[i] = getOrderDomainFromModel(model)
	}

	return orders, nil
}
