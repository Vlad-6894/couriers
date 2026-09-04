package dispetch_postgres_repository

import (
	"context"
	pkg_postgres_pool "couriers/pkg/repository/postgres/pool"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *DispetchRepositoryPostgres) SearchCourier(
	ctx context.Context,
	city string,
) (int, int, error) {
	ctxWithTime, cancel := context.WithTimeout(ctx, r.pool.GetTimeot())
	defer cancel()

	var (
		courierID int
		version   int
	)

	tx, err := r.pool.BeginTx(ctxWithTime, pgx.TxOptions{})
	if err != nil {
		return 0, 0, fmt.Errorf("fail to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	sqlRequestSelect := `
	SELECT id, version FROM app.couriers 
	WHERE city = $1 AND is_free = true
	LIMIT 1
	FOR UPDATE SKIP LOCKED;
	`
	row := tx.QueryRow(ctxWithTime, sqlRequestSelect, city)

	if err := row.Scan(
		&courierID,
		&version,
	); err != nil {
		err = pkg_postgres_pool.MapError(err)
		return 0, 0, fmt.Errorf("fail to scan values: %w", err)
	}

	sqlUpdateRequest := `
	UPDATE app.couriers 
	SET is_free = true, version = version + 1
	WHERE id = $1;
	`
	if _, err := tx.Exec(ctxWithTime, sqlUpdateRequest, courierID); err != nil {
		err = pkg_postgres_pool.MapError(err)
		return 0, 0, fmt.Errorf("fail to update courier: %w", err)
	}

	if err := tx.Commit(ctxWithTime); err != nil {
		return 0, 0, fmt.Errorf("fail to commit: %w", err)
	}

	return courierID, version, nil
}
