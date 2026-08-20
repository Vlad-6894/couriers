package dispetch_redis_repository

import (
	"context"
	dispetch_domains "couriers/internal/Dispetch/core/domains"
	pkg_repository_redis "couriers/pkg/repository/redis"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	cityTTL = 15 * time.Minute
)

func (r *RedisRepository) UpdateCache(
	ctx context.Context,
	couriersByGroupCity map[string]map[int]dispetch_domains.FreeCourierInfo,
) error {
	pipe := r.client.Pipeline()

	delChecks := make(map[string]*redis.IntCmd)
	expireChecks := make(map[string]*redis.BoolCmd)

	for city, couriersMap := range couriersByGroupCity {
		key := fmt.Sprintf("%s:%s", KeyPrefix, city)

		delChecks[city] = pipe.Del(ctx, key)

		for _, courier := range couriersMap {
			value := fmt.Sprintf("%s:%s", strconv.Itoa(courier.ID), strconv.Itoa(courier.Version))

			pipe.SAdd(ctx, key, value)
		}

		expireChecks[city] = pipe.Expire(ctx, key, cityTTL)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		err = pkg_repository_redis.MapError(err)
		return fmt.Errorf("redis pipline error: %w", err)
	}

	for city, cmd := range delChecks {
		if _, err := cmd.Result(); err != nil {
			return fmt.Errorf("fail to delete old courier info from cache for city=%s: %w", city, err)
		}
	}

	for city, cmd := range expireChecks {
		success, err := cmd.Result()
		if err != nil {
			return fmt.Errorf("error expire for city=%s: %w", city, err)
		}

		if !success {
			return fmt.Errorf("TTL was not for city=%s: %w", city, err)
		}
	}

	return nil
}
