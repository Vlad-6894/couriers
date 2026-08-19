package dispetch_redis_repository

import (
	"context"
	pkg_repository_redis "couriers/pkg/repository/redis"
	"fmt"
	"strconv"
	"strings"
)

var (
	KeyPrefix = "couriers:pool"
)

func (r *RedisRepository) SearchCourier(
	ctx context.Context,
	city string,
) (int, int, error) {
	key := fmt.Sprintf("%s:%s", KeyPrefix, city)
	courier, err := r.client.SPop(ctx, key).Result()
	if err != nil {
		err = pkg_repository_redis.MapError(err)
		return 0, 0, fmt.Errorf("fail to get courier from redis: %w", err)
	}

	parts := strings.Split(courier, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("fail to split courier parts")
	}

	courierID, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("fail to convert courierId from string to int: %w", err)
	}

	version, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("fail to convert courierId from string to int: %w", err)
	}

	return courierID, version, nil
}
