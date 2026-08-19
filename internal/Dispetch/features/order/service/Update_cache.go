package dispetch_service

import (
	"context"
	dispetch_domains "couriers/internal/Dispetch/core/domains"
	"fmt"
	"time"
)

var (
	tickerTime = 15 * time.Minute
)

func (s *OrdersDispetchService) UpdateCache(
	ctx context.Context,
) error {
	ticker := time.NewTicker(tickerTime)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			couriersGroupByCity := make(map[string]map[int]dispetch_domains.FreeCourierInfo)

			couriers, err := s.db.GetFreeCouriers(ctx)
			if err != nil {
				return fmt.Errorf("fail to get free couriers from database: %w", err)
			}

			for _, courier := range couriers {
				if _, exists := couriersGroupByCity[courier.City]; !exists {
					cityMap := make(map[int]dispetch_domains.FreeCourierInfo)
					couriersGroupByCity[courier.City] = cityMap
				}

				courierInfo := dispetch_domains.NewFreeCourierInfo(courier.ID, courier.Version)
				couriersGroupByCity[courier.City][courier.ID] = courierInfo

				if err := s.cache.UpdateCache(ctx, couriersGroupByCity); err != nil {
					return fmt.Errorf("fail to update cache from repository: %w", err)
				}
			}
		}
	}
}
