package auth_service

import (
	"context"
	auth_domains "couriers/internal/Auth/core/domains"
	"fmt"
)

func (s *AuthService) RegisterUser(
	ctx context.Context,
	user auth_domains.User,
) (auth_domains.User, error) {
	if err := user.Validate(); err != nil {
		return auth_domains.User{}, fmt.Errorf("validate error: %w", err)
	}

	newUser, err := s.db.RegisterUser(ctx, user)
	if err != nil {
		return auth_domains.User{}, fmt.Errorf("register user repository error: %w", err)
	}

	return newUser, nil
}
