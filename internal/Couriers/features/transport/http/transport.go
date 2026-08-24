package couriers_http_transport

import (
	couriers_transport "couriers/internal/Couriers/features/transport"
	pkg_http_server "couriers/pkg/transport/http/server"
	"net/http"
)

type CouriersHTTPHandler struct {
	service couriers_transport.CouriersService
}

type DispetchedOrderResponseDTO struct {
	OrderID   int `json:"order_id"`
	Version   int `json:"version"`
	CourierID int `json:"courier_id"`
}

func NewCouriersHTTPHandler(
	service couriers_transport.CouriersService,
) *CouriersHTTPHandler {
	return &CouriersHTTPHandler{
		service: service,
	}
}

func (h *CouriersHTTPHandler) Routes() []pkg_http_server.Route {
	routes := []pkg_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/couriers/order",
			Handler: h.HandleGetOrder,
		},

		{
			Method:  http.MethodPatch,
			Path:    "/couriers/order",
			Handler: h.HandleConfirm,
		},
	}

	return routes
}
