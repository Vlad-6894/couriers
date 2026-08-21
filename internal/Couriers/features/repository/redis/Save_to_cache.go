package couriers_redis_repository

import (
	"context"
	courier_domains "couriers/internal/Couriers/core/domains"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	keySets  = "free:couriers:pool"
	orderTTL = 24 * time.Hour
)

func (r *CouriersCache) SaveToCache(
	ctx context.Context,
	order courier_domains.DispetchedOrder,
) error {
	ctxWithTime, cancel := context.WithTimeout(ctx, r.client.GetTimeout())
	cancel()

	value := fmt.Sprintf("%s:%s:%s", strconv.Itoa(order.CourierID), strconv.Itoa(order.OrderID), strconv.Itoa(order.Version))

	if err := r.client.ZAdd(ctxWithTime, keySets, redis.Z{
		Score:  0,
		Member: value,
	}).Err(); err != nil {
		return fmt.Errorf("fail to save cache: %w", err)
	}

	return nil
}
