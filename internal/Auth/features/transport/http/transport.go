package auth_http_transport

import (
	"context"
	auth_domains "couriers/internal/Auth/core/domains"
	pkg_http_server "couriers/pkg/transport/http/server"
	"net/http"
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

	LoginUser(ctx context.Context, login string, password string) (string, error)

	LoginCourier(ctx context.Context, login string, password string) (string, error)
}

type RegisterRequestDTO struct {
	Login    string `json:"Login"      example:"123456789"`
	Password string `json:"Password"   example:"123456789"`
	City     string `json:"City"       example:"Moscow"`
}

type RegisterUserResponseDTO struct {
	ID       int    `json:"User_ID"        example:"1"`
	Version  int    `json:"User_version"   example:"1"`
	Login    string `json:"Login"          example:"123456789"`
	Password string `json:"Password"       example:"123456789"`
	City     string `json:"City"           example:"Moscow"`
}

type LoginRequestDTO struct {
	Login    string `json:"Login"          example:"123456789"`
	Password string `json:"Password"       example:"123456789"`
}

type CourierResponseDTO struct {
	ID             int    `json:"User_ID"                example:"1"`
	Version        int    `json:"Version"                example:"1"`
	Login          string `json:"Login"                  example:"123456789"`
	Password       string `json:"Password"               example:"123456789"`
	City           string `json:"City"                   example:"Moscow"`
	OrdersComplete int    `json:"orders_complete"        example:"1"`
	IsFree         bool   `json:"is_free"                example:"true"`
}

func NewAuthHTTPHandler(authService AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}

func (h *AuthHTTPHandler) UserRoutes() []pkg_http_server.Route {
	return []pkg_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "auth/register",
			Handler: h.HandleRegisterUser,
		},

		{
			Method:  http.MethodPost,
			Path:    "auth/login",
			Handler: h.HandleLoginUser,
		},
	}
}

func (h *AuthHTTPHandler) CourierRoutes() []pkg_http_server.Route {
	return []pkg_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "auth/register",
			Handler: h.HandleRegisterCourier,
		},

		{
			Method:  http.MethodPost,
			Path:    "auth/login",
			Handler: h.HandleLoginCourier,
		},
	}
}
