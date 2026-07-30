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

	RegisterCourier(
		ctx context.Context,
		courier auth_domains.Courier,
	) (auth_domains.Courier, error)
}

type RegisterRequestDTO struct {
	Login    string `json:"Login"`
	Password string `json:"Password"`
	City     string `json:"City"`
}

type CourierResponseDTO struct {
	ID             int    `json:"User_ID"`
	Version        int    `json:"Version"`
	Login          string `json:"Login"`
	Password       string `json:"Password"`
	City           string `json:"City"`
	OrdersComplete int    `json:"orders_complete"`
	IsFree         bool   `json:"is_free"`
}

func NewAuthHTTPHandler(authService AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}
