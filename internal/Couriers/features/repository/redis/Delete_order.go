package couriers_redis_repository

import (
	"context"
	"fmt"
	"strconv"
)

func (r *CouriersCache) DeleteOrder(
	ctx context.Context,
	personID int,
) error {
	minRange := fmt.Sprintf("[courier:%s#", strconv.Itoa(personID))
	maxRange := fmt.Sprintf("[courier:%s#\xff", strconv.Itoa(personID))

	_, err := r.client.ZRemRangeByLex(ctx, keySets, minRange, maxRange).Result()
	if err != nil {
		return fmt.Errorf("fail delete:: %w", err)
	}

	return nil
}
