package auth_http_transport

import (
	auth_domains "couriers/internal/Auth/core/domains"
	pkg_logger "couriers/pkg/logger"
	pkg_http_response "couriers/pkg/transport/http/response"
	"encoding/json"
	"net/http"
)

type CreateCourierResponseDTO CourierResponseDTO

func (h *AuthHTTPHandler) HandleRegisterCourier(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := pkg_logger.FromContext(ctx)
	responseHandler := pkg_http_response.NewHTTPResponseHandler(log, w)

	log.Info("Start Register Courier")

	var request RegisterRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responseHandler.ErrorResponse(err, "decode json error: ")
		return
	}

	courier := regCourierDtoToUserDomain(request)

	courier, err := h.authService.RegisterCourier(ctx, courier)
	if err != nil {
		responseHandler.ErrorResponse(err, "error Register Courier service error: ")
		return
	}

	response := CreateCourierResponseDTO(regCourierDtoFromUserDomain(courier))

	responseHandler.ToJSONResponse(response, http.StatusCreated)

	log.Info("Finish Register Courier")
}

func regCourierDtoToUserDomain(dto RegisterRequestDTO) auth_domains.Courier {
	courier := auth_domains.NewRegCourier(dto.Login, dto.Password, dto.City)
	return courier
}

func regCourierDtoFromUserDomain(courier auth_domains.Courier) CourierResponseDTO {
	return CourierResponseDTO{
		ID:             courier.ID,
		Version:        courier.Version,
		Login:          courier.Login,
		Password:       courier.Password,
		City:           courier.City,
		OrdersComplete: courier.OrdersComplete,
		IsFree:         courier.IsFree,
	}
}
