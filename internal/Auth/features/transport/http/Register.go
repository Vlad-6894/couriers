package auth_http_transport

import (
	auth_domains "couriers/internal/Auth/core/domains"
	pkg_logger "couriers/pkg/logger"
	pkg_http_response "couriers/pkg/transport/http/response"
	"encoding/json"
	"net/http"
)

type RegisterRequestDTO struct {
	Login    string `json:"Login"`
	Password string `json:"Password"`
	City     string `json:"City"`
}

func (h *AuthHTTPHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := pkg_logger.FromContext(ctx)
	responseHandler := pkg_http_response.NewHTTPResponseHandler(log, w)

	log.Info("Start Register person")

	var request RegisterRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responseHandler.ErrorResponse(err, "fail to decode from json!")
		return
	}

	user := RegDtoToUserDomain(request)

	if err := h.authService.Register(user); err != nil {
		responseHandler.ErrorResponse(err, "fail to register user!")
		return
	}

	w.WriteHeader(http.StatusCreated)

	log.Info("Finish Register person")
}

func RegDtoToUserDomain(dto RegisterRequestDTO) auth_domains.User {
	user := auth_domains.NewRegUser(dto.Login, dto.Password, dto.City)
	return user
}
