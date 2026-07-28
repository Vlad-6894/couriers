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
}

func NewAuthService(
	db AuthDatabaseRepository,
) *AuthService {
	return &AuthService{
		db: db,
	}
}
