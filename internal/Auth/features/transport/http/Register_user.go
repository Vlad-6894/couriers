package auth_http_transport

import (
	auth_domains "couriers/internal/Auth/core/domains"
	pkg_logger "couriers/pkg/logger"
	pkg_http_response "couriers/pkg/transport/http/response"
	"encoding/json"
	"net/http"
)

type CreateUserDTO RegisterUserResponseDTO

func (h *AuthHTTPHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := pkg_logger.FromContext(ctx)
	responseHandler := pkg_http_response.NewHTTPResponseHandler(log, w)

	log.Info("Start Register User")

	var request RegisterRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responseHandler.ErrorResponse(err, "fail to decode from json!")
		return
	}

	user := regUserDtoToUserDomain(request)

	user, err := h.authService.RegisterUser(ctx, user)
	if err != nil {
		responseHandler.ErrorResponse(err, "fail to register user!")
		return
	}

	response := CreateUserDTO(regUserDtoFromUserDomain(user))

	responseHandler.ToJSONResponse(response, http.StatusCreated)

	log.Info("Finish Register User")
}

func regUserDtoToUserDomain(dto RegisterRequestDTO) auth_domains.User {
	user := auth_domains.NewRegUser(dto.Login, dto.Password, dto.City)
	return user
}

func regUserDtoFromUserDomain(user auth_domains.User) RegisterUserResponseDTO {
	return RegisterUserResponseDTO{
		ID:       user.ID,
		Version:  user.Version,
		Login:    user.Login,
		Password: user.Password,
		City:     user.City,
	}
}
