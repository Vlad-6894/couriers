package dispetch_confirm_repository_postgres

import pkg_postgres_pool "couriers/pkg/repository/postgres/pool"

type ConfirmPostgresDatabase struct {
	pool pkg_postgres_pool.Pool
}

func NewConfirmPostgresDatabase(
	pool pkg_postgres_pool.Pool,
) *ConfirmPostgresDatabase {
	return &ConfirmPostgresDatabase{
		pool: pool,
	}
}
