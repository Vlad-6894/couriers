package dispetch_postgres_repository

import pkg_postgres_pool "couriers/pkg/repository/postgres/pool"

type DispetchRepositoryPostgres struct {
	pool pkg_postgres_pool.Pool
}

func NewDispetchRepositoryPostgres(
	pool pkg_postgres_pool.Pool,
) *DispetchRepositoryPostgres {
	return &DispetchRepositoryPostgres{
		pool: pool,
	}
}
