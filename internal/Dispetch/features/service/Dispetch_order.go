package dispetch_service

import (
	"context"
	dispetch_domains "couriers/internal/Dispetch/core/domains"
	pkg_repository_redis "couriers/pkg/repository/redis"
	"errors"
	"fmt"
)

func (s *OrdersDispetchService) DispetchOrder(
	ctx context.Context,
	order dispetch_domains.Order,
) error {
	courierID, version, err := s.cache.SearchCourier(ctx, order.City)
	if err != nil {
		if errors.Is(err, pkg_repository_redis.ErrEmpty) {
			courierID, version, err := s.db.SearchCourier(ctx, order.City)
			if err != nil {
				return fmt.Errorf("fail to get courier id: %w", err)
			}
			if err := s.broker.SendDispetchedOrder(ctx, order.ID, version, courierID, order.City); err != nil {
				return fmt.Errorf("fail to send order to kafka: %w", err)
			}
		}

		return fmt.Errorf("fail to cache: %w", err)

	}

	if err := s.db.DoBusy(ctx, courierID, version); err != nil {
		return fmt.Errorf("fail to do courier busy: %w", err)
	}

	if err := s.broker.SendDispetchedOrder(ctx, order.ID, version, courierID, order.City); err != nil {
		return fmt.Errorf("fail to send order to kafka: %w", err)
	}

	return nil
}
