package dispetch_redis_repository

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

var (
	TTL = 5 * time.Minute
)

func (r *RedisRepository) CheckUnique(
	ctx context.Context,
	orderId int,
) (bool, error) {
	key := fmt.Sprintf("idempotency:order:%s", strconv.Itoa(orderId))

	succes, err := r.client.SetNX(ctx, key, orderId, TTL).Result()
	if err != nil {
		return succes, fmt.Errorf("fail to setNX redis: %w", err)
	}

	return succes, nil
}
