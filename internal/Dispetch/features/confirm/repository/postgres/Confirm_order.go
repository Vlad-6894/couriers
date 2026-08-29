package dispetch_confirm_repository_postgres

import (
	"context"
	dispetch_domains "couriers/internal/Dispetch/core/domains"
	pkg_postgres_pool "couriers/pkg/repository/postgres/pool"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *ConfirmPostgresDatabase) ConfirmOrder(
	ctx context.Context,
	confirm dispetch_domains.Confirm,
) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, r.pool.GetTimeot())
	defer cancel()

	tx, err := r.pool.BeginTx(ctxWithTimeout, pgx.TxOptions{})
	if err != nil {
		err = pkg_postgres_pool.MapError(err)
		return fmt.Errorf("fail to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	sqlUpdateOrders := `
	UPDATE app.orders 
	SET version = version + 1, is_complete = true, courier_id = $1
	WHERE id = $2;
	`

	if _, err := r.pool.Exec(
		ctxWithTimeout,
		sqlUpdateOrders,
		confirm.CourierID,
		confirm.OrderID,
	); err != nil {
		err = pkg_postgres_pool.MapError(err)
		return fmt.Errorf("fail to update orders table: %w", err)
	}

	sqlUpdateCouriers := `
	UPDATE app.couriers
	SET version = version + 1, is_free = true, orders_complete = orders_complete + 1
	WHERE id = $1;
	`

	if _, err := r.pool.Exec(
		ctxWithTimeout,
		sqlUpdateCouriers,
		confirm.CourierID,
	); err != nil {
		err = pkg_postgres_pool.MapError(err)
		return fmt.Errorf("fail to update couriers table: %w", err)
	}

	if err := tx.Commit(ctxWithTimeout); err != nil {
		return fmt.Errorf("fail to commit: %w", err)
	}

	return nil
}
