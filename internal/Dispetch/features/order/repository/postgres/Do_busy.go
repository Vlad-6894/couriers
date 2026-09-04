package dispetch_postgres_repository

import (
	"context"
	pkg_postgres_pool "couriers/pkg/repository/postgres/pool"
	"fmt"
)

func (r *DispetchRepositoryPostgres) DoBusy(
	ctx context.Context,
	courierID int,
	version int,
) error {
	ctxWithTime, cancel := context.WithTimeout(ctx, r.pool.GetTimeot())
	defer cancel()

	sqlRequest := `
	UPDATE app.couriers 
	SET is_free = false, version = version + 1
	WHERE id = $1 AND version = $2;
	`

	if _, err := r.pool.Exec(ctxWithTime, sqlRequest, courierID, version); err != nil {
		err = pkg_postgres_pool.MapError(err)
		return fmt.Errorf("fail to exec: %w", err)
	}

	return nil
}
