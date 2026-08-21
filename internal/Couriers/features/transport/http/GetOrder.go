package couriers_http_transport

import (
	courier_domains "couriers/internal/Couriers/core/domains"
	pkg_jwt "couriers/pkg/jwt"
	pkg_logger "couriers/pkg/logger"
	pkg_http_response "couriers/pkg/transport/http/response"
	"net/http"
)

func (h *CouriersHTTPHandler) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := pkg_logger.FromContext(ctx)
	responseHandler := pkg_http_response.NewHTTPResponseHandler(log, w)
	personID := pkg_jwt.PersonIDFromContext(ctx)

	order, err := h.service.GetOrder(ctx, personID)
	if err != nil {
		responseHandler.ErrorResponse(err, "fail to get order!")
		return
	}

	response := orderDtoFromDomain(order)

	responseHandler.ToJSONResponse(response, http.StatusOK)
}

func orderDtoFromDomain(order courier_domains.DispetchedOrder) DispetchedOrderResponseDTO {
	dto := DispetchedOrderResponseDTO{
		OrderID:   order.OrderID,
		Version:   order.Version,
		CourierID: order.CourierID,
	}

	return dto
}
