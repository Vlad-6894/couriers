package dispetch_postgres_repository

import (
	dispetch_domains "couriers/internal/Dispetch/core/domains"
	pkg_postgres_pool "couriers/pkg/repository/postgres/pool"
)

type DispetchRepositoryPostgres struct {
	pool pkg_postgres_pool.Pool
}

type CourierInfoModel struct {
	ID      int
	Version int
	City    string
}

func NewDispetchRepositoryPostgres(
	pool pkg_postgres_pool.Pool,
) *DispetchRepositoryPostgres {
	return &DispetchRepositoryPostgres{
		pool: pool,
	}
}

func domainInfoFromModel(models []CourierInfoModel) []dispetch_domains.CourierInfo {
	info := make([]dispetch_domains.CourierInfo, len(models))
	for i, v := range models {
		info[i] = dispetch_domains.NewCourierInfo(v.ID, v.Version, v.City)
	}

	return info
}
