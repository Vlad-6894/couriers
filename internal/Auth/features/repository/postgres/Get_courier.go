package auth_repository_postgres

import (
	"context"
	auth_domains "couriers/internal/Auth/core/domains"
	"fmt"
)

func (r *AuthRepository) GetCourier(ctx context.Context, login string) (auth_domains.Courier, error) {
	ctxWithTime, cancel := context.WithTimeout(ctx, r.pool.GetTimeot())
	defer cancel()

	sqlRequest := `
	SELECT id, version, login, password, city, orders_complete, is_free FROM app.couriers
	WHERE login = $1;
	`

	row := r.pool.QueryRow(ctxWithTime, sqlRequest, login)

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
		return auth_domains.Courier{}, fmt.Errorf("error get courier from database %w", err)
	}

	courier := NewCourierDomainFromModel(courierModel)

	return courier, nil
}
