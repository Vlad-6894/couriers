package couriers_redis_repository

import pkg_repository_redis "couriers/pkg/repository/redis"

type CouriersCache struct {
	client pkg_repository_redis.ClientCacheRedis
}

func NewCouriersCache(
	client pkg_repository_redis.ClientCacheRedis,
) *CouriersCache {
	return &CouriersCache{
		client: client,
	}
}
