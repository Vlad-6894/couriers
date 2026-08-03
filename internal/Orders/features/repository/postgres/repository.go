package orders_repository_posgres

import pkg_postgres_pool "couriers/pkg/repository/postgres/pool"

type DatabasePostgres struct {
	pool pkg_postgres_pool.Pool
}

func NewDatabasePostgres(pool pkg_postgres_pool.Pool) *DatabasePostgres {
	return &DatabasePostgres{
		pool: pool,
	}
}
