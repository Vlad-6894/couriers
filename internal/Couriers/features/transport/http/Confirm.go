package couriers_http_transport

import (
	pkg_jwt "couriers/pkg/jwt"
	pkg_logger "couriers/pkg/logger"
	pkg_http_response "couriers/pkg/transport/http/response"
	"net/http"
)

func (h *CouriersHTTPHandler) HandleConfirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := pkg_logger.FromContext(ctx)
	responseHandler := pkg_http_response.NewHTTPResponseHandler(log, w)
	personID := pkg_jwt.PersonIDFromContext(ctx)

	if err := h.service.Confirm(ctx, personID); err != nil {
		responseHandler.ErrorResponse(err, "fail to confirm order from service!")
		return
	}

	responseHandler.ToJSONResponse("Confirm was success!", http.StatusNoContent)
}
