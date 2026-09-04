package auth_service

import (
	"context"
	auth_domains "couriers/internal/Auth/core/domains"
	"fmt"
)

func (s *AuthService) RegisterCourier(
	ctx context.Context,
	courier auth_domains.Courier,
) (auth_domains.Courier, error) {
	if err := courier.Validate(); err != nil {
		return auth_domains.Courier{}, fmt.Errorf("Register Courier validate error: %w", err)
	}

	newCourier, err := s.db.RegisterCourier(ctx, courier)
	if err != nil {
		return auth_domains.Courier{}, fmt.Errorf("Register Courier repository error: %w", err)
	}

	return newCourier, nil
}
