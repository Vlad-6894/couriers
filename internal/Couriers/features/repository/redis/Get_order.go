package couriers_redis_repository

import (
	"context"
	courier_domains "couriers/internal/Couriers/core/domains"
	"fmt"
	"strconv"
	"strings"
)

func (r *CouriersCache) GetOrder(
	ctx context.Context,
	personID int,
) (courier_domains.DispetchedOrder, error) {
	ctxWithTime, cancel := context.WithTimeout(ctx, r.client.GetTimeout())
	defer cancel()

	prefix := fmt.Sprintf("%s:*", strconv.Itoa(personID))

	courier, _, err := r.client.ZScan(ctxWithTime, keySets, 0, prefix, 100).Result()
	if err != nil {
		return courier_domains.DispetchedOrder{}, fmt.Errorf("fail to get order: %w", err)
	}

	courierValue := courier[0]

	sliceValues := strings.Split(courierValue, ":")
	courierIDStr := sliceValues[0]
	orderIdStr := sliceValues[1]
	versionStr := sliceValues[2]

	courierID, err := strconv.Atoi(courierIDStr)
	if err != nil {
		return courier_domains.DispetchedOrder{}, fmt.Errorf("fail to strconv: %w", err)
	}

	orderID, err := strconv.Atoi(orderIdStr)
	if err != nil {
		return courier_domains.DispetchedOrder{}, fmt.Errorf("fail to strconv: %w", err)
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return courier_domains.DispetchedOrder{}, fmt.Errorf("fail to strconv: %w", err)
	}

	order := courier_domains.NewDispetchedOrder(orderID, version, courierID)

	return order, nil
}
