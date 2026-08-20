package dispetch_confirm_service

import (
	"context"
	dispetch_domains "couriers/internal/Dispetch/core/domains"
	"fmt"
)

func (s *DispetchConfirmService) ConfirmOrder(
	ctx context.Context,
	confirm dispetch_domains.Confirm,
) error {
	if err := s.db.ConfirmOrder(ctx, confirm); err != nil {
		return fmt.Errorf("fail to confirm from repository: %w", err)
	}

	return nil
}
