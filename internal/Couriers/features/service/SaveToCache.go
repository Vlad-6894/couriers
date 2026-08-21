package couriers_service

import (
	"context"
	courier_domains "couriers/internal/Couriers/core/domains"
	"fmt"
)

func (s *CouriersService) SaveToCache(
	ctx context.Context,
	order courier_domains.DispetchedOrder,
) error {
	if err := s.cache.SaveToCache(ctx, order); err != nil {
		return fmt.Errorf("fail to save cache from repository: %w", err)
	}

	return nil
}
