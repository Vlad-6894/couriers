package orders_http_transport

import (
	"context"
	orders_domains "couriers/internal/Orders/core/domains"
	pkg_http_server "couriers/pkg/transport/http/server"
	"net/http"
)

type OrdersHTTPHandler struct {
	ordersService OrdersService
}

type OrdersService interface {
	CreateOrder(
		ctx context.Context,
		order orders_domains.Order,
	) (orders_domains.Order, error)

	GetOrders(
		ctx context.Context,
		personID int,
		limit *int,
		offset *int,
	) ([]orders_domains.GetOrder, error)
}

type OrderResponseDTO struct {
	ID         int    `json:"order_id"`
	Version    int    `json:"order_version"`
	Name       string `json:"order_name"`
	Price      int    `json:"order_price"`
	IsComplete bool   `json:"is_complete"`
	UserID     int    `json:"user_id"`
	CourierID  *int   `json:"courier_id"`
}

type GetOrderResponseDTO struct {
	ID           int     `json:"order_id"`
	Version      int     `json:"order_version"`
	Name         string  `json:"order_name"`
	Price        int     `json:"order_price"`
	IsComplete   bool    `json:"is_complete"`
	UserID       int     `json:"user_id"`
	CourierID    *int    `json:"courier_id"`
	CourierLogin *string `json:"courier_login"`
}

func NewOrdersHTTPHandler(ordersService OrdersService) *OrdersHTTPHandler {
	return &OrdersHTTPHandler{
		ordersService: ordersService,
	}
}

func orderDtoFromOrderDomain(order orders_domains.Order) OrderResponseDTO {
	return OrderResponseDTO{
		ID:         order.ID,
		Version:    order.Version,
		Name:       order.Name,
		Price:      order.Price,
		IsComplete: order.IsComplete,
		UserID:     order.UserID,
		CourierID:  order.CourierID,
	}
}

func NewGetOrderResponseDTO(
	id int,
	version int,
	name string,
	price int,
	isComplete bool,
	userID int,
	courierID *int,
	courierLogin *string,
) GetOrderResponseDTO {
	return GetOrderResponseDTO{
		ID:           id,
		Version:      version,
		Name:         name,
		Price:        price,
		IsComplete:   isComplete,
		UserID:       userID,
		CourierID:    courierID,
		CourierLogin: courierLogin,
	}
}

func (h *OrdersHTTPHandler) Routes() []pkg_http_server.Route {
	routes := []pkg_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/order",
			Handler: h.HandleCreateOrder,
		},

		{
			Method:  http.MethodGet,
			Path:    "/order",
			Handler: h.HandleGetOrders,
		},
	}

	return routes
}
