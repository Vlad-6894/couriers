package orders_service

import (
	"context"
	orders_domains "couriers/internal/Orders/core/domains"
	"fmt"
)

func (s *OrdersService) CreateOrder(
	ctx context.Context,
	order orders_domains.Order,
) (orders_domains.Order, error) {
	if err := order.Validate(); err != nil {
		return orders_domains.Order{}, fmt.Errorf("fail validate order: %w", err)
	}

	order, err := s.db.CreateOrder(ctx, order)
	if err != nil {
		return orders_domains.Order{}, fmt.Errorf("fail to save order to database: %w", err)
	}

	if err := s.broker.SendOrder(ctx, order); err != nil {
		return orders_domains.Order{}, fmt.Errorf("fail to send order to kafka: %w", err)
	}

	return order, nil
}
