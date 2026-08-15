package dispetch_redis_repository

import pkg_repository_redis "couriers/pkg/repository/redis"

type RedisRepository struct {
	client pkg_repository_redis.ClientCacheRedis
}

func NewRedisRepository(
	client pkg_repository_redis.ClientCacheRedis,
) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}
