package auth_http_transport

import (
	pkg_logger "couriers/pkg/logger"
	pkg_http_response "couriers/pkg/transport/http/response"
	"encoding/json"
	"net/http"
)

// LoginUser godoc
// @Summary Аутентифицировать пользователя
// @Description Аутентифицировать пользователя в системе
// @Tags users
// @Accept json
// @Produce plain
// @Param request body LoginRequestDTO true "тело запроса"
// @Failure 400 {object} pkg_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} pkg_http_response.ErrorResponse "Internal server error"
// @Router /users/auth/login [post]
func (h *AuthHTTPHandler) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := pkg_logger.FromContext(ctx)
	responseHandler := pkg_http_response.NewHTTPResponseHandler(log, w)

	log.Info("Start Login user")

	var request LoginRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responseHandler.ErrorResponse(err, "fail to decode login request error!")
		return
	}

	token, err := h.authService.LoginUser(ctx, request.Login, request.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "login user error!")
		return
	}

	responseHandler.ToJSONResponse(token, http.StatusOK)

	log.Info("Finish Login user")
}
