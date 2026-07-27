package auth_http_transport

import auth_domains "couriers/internal/Auth/core/domains"

type AuthHTTPHandler struct {
	authService AuthService
}

type AuthService interface {
	Register(auth_domains.User) error
}

func NewAuthHTTPHandler(authService AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}
