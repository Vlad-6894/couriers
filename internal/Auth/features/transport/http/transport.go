package auth_http_transport

import (
	"context"
	auth_domains "couriers/internal/Auth/core/domains"
)

type AuthHTTPHandler struct {
	authService AuthService
}

type AuthService interface {
	RegisterUser(
		ctx context.Context,
		user auth_domains.User,
	) (auth_domains.User, error)
}

type RegisterRequestDTO struct {
	Login    string `json:"Login"`
	Password string `json:"Password"`
	City     string `json:"City"`
}

func NewAuthHTTPHandler(authService AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}
