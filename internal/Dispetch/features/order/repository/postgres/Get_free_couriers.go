package dispetch_postgres_repository

import (
	"context"
	dispetch_domains "couriers/internal/Dispetch/core/domains"
	pkg_postgres_pool "couriers/pkg/repository/postgres/pool"
	"fmt"
)

func (r *DispetchRepositoryPostgres) GetFreeCouriers(
	ctx context.Context,
) ([]dispetch_domains.CourierInfo, error) {
	ctxWithTime, cancel := context.WithTimeout(ctx, r.pool.GetTimeot())
	cancel()

	sqlRequest := `
	SELECT id, version, city FROM app.couriers 
	WHERE is_free = true;
	`

	rows, err := r.pool.Query(ctxWithTime, sqlRequest)
	if err != nil {
		err = pkg_postgres_pool.MapError(err)
		return nil, fmt.Errorf("fail to get rows from database: %w", err)
	}
	rows.Close()

	models := make([]CourierInfoModel, 0)

	for rows.Next() {
		var model CourierInfoModel
		if err := rows.Scan(
			&model.ID,
			&model.Version,
			&model.City,
		); err != nil {
			err = pkg_postgres_pool.MapError(err)
			return nil, fmt.Errorf("fail to scan: %w", err)
		}

		models = append(models, model)
	}

	info := domainInfoFromModel(models)

	return info, nil
}
