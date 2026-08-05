package orders_service

import (
	"context"
	orders_domains "couriers/internal/Orders/core/domains"
	pkg_errors "couriers/pkg/errors"
	"fmt"
)

func (s *OrdersService) GetOrders(
	ctx context.Context,
	personID int,
	limit *int,
	offset *int,
) ([]orders_domains.GetOrder, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be positive or 0: %w", pkg_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be positive or 0: %w", pkg_errors.ErrInvalidArgument)
	}

	orders, err := s.db.GetOrders(
		ctx,
		personID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get orders from database error: %w", err)
	}

	return orders, nil
}
