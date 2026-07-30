package auth_repository_postgres

import (
	"context"
	auth_domains "couriers/internal/Auth/core/domains"
	"fmt"
)

func (r *AuthRepository) RegisterCourier(
	ctx context.Context,
	courier auth_domains.Courier,
) (auth_domains.Courier, error) {
	ctxWithTime, cancel := context.WithTimeout(ctx, r.pool.GetTimeot())
	defer cancel()

	sqlRequest := `
	INSERT INTO app.Couriers (login, password, city, orders_complete, is_free)
	VALUES($1,$2,$3,$4,$5)
	RETURNING id, version, login, password, city, orders_complete, is_free;
	`

	row := r.pool.QueryRow(
		ctxWithTime,
		sqlRequest,
		courier.Login,
		courier.Password,
		courier.City,
		courier.OrdersComplete,
		courier.IsFree,
	)

	var courierModel CourierModel

	if err := row.Scan(
		&courierModel.ID,
		&courierModel.Version,
		&courierModel.Login,
		&courierModel.Password,
		&courierModel.City,
		&courierModel.OrdersComplete,
		&courierModel.IsFree,
	); err != nil {
		return auth_domains.Courier{}, fmt.Errorf("error scan: %w", err)
	}

	newCourier := NewCourierDomainFromModel(courierModel)

	return newCourier, nil
}
