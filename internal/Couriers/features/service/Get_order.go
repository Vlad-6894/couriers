package couriers_service

import (
	"context"
	courier_domains "couriers/internal/Couriers/core/domains"
	"fmt"
)

func (s *CouriersService) GetOrder(
	ctx context.Context,
	personID int,
) (courier_domains.DispetchedOrder, error) {
	order, err := s.cache.GetOrder(ctx, personID)
	if err != nil {
		return courier_domains.DispetchedOrder{}, fmt.Errorf("fail to get from cache: %w", err)
	}

	return order, nil
}
