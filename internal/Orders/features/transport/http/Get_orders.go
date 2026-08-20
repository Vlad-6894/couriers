package orders_http_transport

import (
	orders_domains "couriers/internal/Orders/core/domains"
	pkg_jwt "couriers/pkg/jwt"
	pkg_logger "couriers/pkg/logger"
	pkg_http_response "couriers/pkg/transport/http/response"
	pkg_http_utils "couriers/pkg/transport/http/utils"
	"net/http"
)

type GetOrdersResponseDTO []GetOrderResponseDTO

func (h *OrdersHTTPHandler) HandleGetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := pkg_logger.FromContext(ctx)
	personID := pkg_jwt.PersonIDFromContext(ctx)
	responseHandler := pkg_http_response.NewHTTPResponseHandler(log, w)

	limit, offset, err := pkg_http_utils.GetLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "fail get limit and offset!")
		return
	}

	orders, err := h.ordersService.GetOrders(
		ctx,
		personID,
		limit,
		offset,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "fail to get orders!")
		return
	}

	response := GetOrdersResponseDTO(ordersDtoFromDomains(orders))

	responseHandler.ToJSONResponse(response, http.StatusOK)
}

func ordersDtoFromDomains(orders []orders_domains.GetOrder) []GetOrderResponseDTO {
	ordersdto := make([]GetOrderResponseDTO, len(orders))

	for i, order := range orders {
		dto := NewGetOrderResponseDTO(
			order.ID,
			order.Version,
			order.Name,
			order.Price,
			order.IsComplete,
			order.City,
			order.UserLogin,
			order.UserID,
			order.CourierID,
			order.CourierLogin,
		)

		ordersdto[i] = dto
	}

	return ordersdto
}
