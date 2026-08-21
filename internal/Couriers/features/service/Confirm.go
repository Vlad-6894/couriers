package couriers_service

import (
	"context"
	"fmt"
)

func (s *CouriersService) Confirm(
	ctx context.Context,
	personID int,
) error {
	orderInfo, err := s.cache.GetOrder(ctx, personID)
	if err != nil {
		return fmt.Errorf("fail to get order from cache: %w", err)
	}

	if err := s.cache.DeleteOrder(ctx, personID); err != nil {
		return fmt.Errorf("fail to delete order from repository: %w", err)
	}

	if err := s.broker.SendConfirm(ctx, orderInfo); err != nil {
		return fmt.Errorf("fail to send confirm: %w", err)
	}

	return nil
}
