package pkg_http_response

import (
	pkg_errors "couriers/pkg/errors"
	pkg_logger "couriers/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *pkg_logger.Logger
	w   http.ResponseWriter
}

func NewHTTPResponseHandler(log *pkg_logger.Logger, w http.ResponseWriter) HTTPResponseHandler {
	return HTTPResponseHandler{
		log: log,
		w:   w,
	}
}

func (h *HTTPResponseHandler) ToJSONResponse(responseBody any, statusCode int) {
	h.w.WriteHeader(statusCode)

	if err := json.NewEncoder(h.w).Encode(responseBody); err != nil {
		h.log.Error("faile to encode: ", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) ErrorResponse(err error, message string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, pkg_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	case errors.Is(err, pkg_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, pkg_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(message, zap.Error(err))

	h.errorResponse(statusCode, message, err)

}

func (h *HTTPResponseHandler) PanicResponse(message string, p any) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("panic! %v", p)

	h.log.Error(message, zap.Error(err))
	h.errorResponse(statusCode, message, err)
}

func (h *HTTPResponseHandler) errorResponse(statusCode int, message string, err error) {
	response := ErrorResponse{
		Error:   err.Error(),
		Message: message,
	}

	h.w.WriteHeader(statusCode)

	h.ToJSONResponse(response, statusCode)
}
