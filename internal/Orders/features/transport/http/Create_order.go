package orders_http_transport

import (
	orders_domains "couriers/internal/Orders/core/domains"
	pkg_jwt "couriers/pkg/jwt"
	pkg_logger "couriers/pkg/logger"
	pkg_http_response "couriers/pkg/transport/http/response"
	"encoding/json"
	"net/http"
)

type CreateOrderRequestDTO struct {
	Name  string `json:"order_name"`
	Price int    `json:"order_price"`
}

type CreateOrderResponseDTO OrderResponseDTO

func (h *OrdersHTTPHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := pkg_logger.FromContext(ctx)
	responseHandler := pkg_http_response.NewHTTPResponseHandler(log, w)

	var request CreateOrderRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responseHandler.ErrorResponse(err, "Fail to decode order DTO!")
		return
	}

	personID := pkg_jwt.PersonIDFromContext(ctx)

	order := orderDomainFromCreateOrderDto(request, personID)

	order, err := h.ordersService.CreateOrder(ctx, order)

	if err != nil {
		responseHandler.ErrorResponse(err, "fail to create order!")
		return
	}

	response := CreateOrderResponseDTO(orderDtoFromOrderDomain(order))

	responseHandler.ToJSONResponse(response, http.StatusCreated)
}

func orderDomainFromCreateOrderDto(dto CreateOrderRequestDTO, personID int) orders_domains.Order {
	order := orders_domains.NewUnitializedOrder(dto.Name, dto.Price, personID)
	return order
}
