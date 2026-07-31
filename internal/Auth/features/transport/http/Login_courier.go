package auth_http_transport

import (
	pkg_logger "couriers/pkg/logger"
	pkg_http_response "couriers/pkg/transport/http/response"
	"encoding/json"
	"net/http"
)

func (h *AuthHTTPHandler) HandleLoginCourier(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := pkg_logger.FromContext(ctx)
	responseHandler := pkg_http_response.NewHTTPResponseHandler(log, w)

	log.Info("start Login Courier!")

	var request LoginRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responseHandler.ErrorResponse(err, "fail to decode login request error!")
		return
	}

	token, err := h.authService.LoginCourier(ctx, request.Login, request.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "login user error!")
		return
	}

	responseHandler.ToJSONResponse(token, http.StatusOK)

	log.Info("Finish Login courier!")
}
