package auth_repository_postgres

import pkg_postgres_pool "couriers/pkg/repository/postgres/pool"

type AuthRepository struct {
	pool pkg_postgres_pool.Pool
}

func NewAuthRepository(pool pkg_postgres_pool.Pool) *AuthRepository {
	return &AuthRepository{
		pool: pool,
	}
}
