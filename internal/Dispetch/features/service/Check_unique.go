package dispetch_service

import (
	"context"
	"fmt"
)

func (s *OrdersDispetchService) CheckUnique(
	ctx context.Context,
	orderId int,
) (bool, error) {
	isUnique, err := s.cache.CheckUnique(ctx, orderId)
	if err != nil {
		return isUnique, fmt.Errorf("fail to check unique from cache: %w", err)
	}

	return isUnique, err
}
