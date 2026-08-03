package orders_http_transport

import (
	"context"
	orders_domains "couriers/internal/Orders/core/domains"
)

type OrdersHTTPHandler struct {
	ordersService OrdersService
}

type OrdersService interface {
	CreateOrder(
		ctx context.Context,
		order orders_domains.Order,
	) (orders_domains.Order, error)
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
