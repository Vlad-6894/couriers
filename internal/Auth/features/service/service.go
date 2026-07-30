package auth_service

import (
	"context"
	auth_domains "couriers/internal/Auth/core/domains"
)

type AuthService struct {
	db AuthDatabaseRepository
}

type AuthDatabaseRepository interface {
	RegisterUser(ctx context.Context, user auth_domains.User) (auth_domains.User, error)
	RegisterCourier(ctx context.Context, courier auth_domains.Courier) (auth_domains.Courier, error)
}

func NewAuthService(
	db AuthDatabaseRepository,
) *AuthService {
	return &AuthService{
		db: db,
	}
}
