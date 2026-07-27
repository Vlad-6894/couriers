package auth_core_http_middleware

import (
	pkg_logger "couriers/pkg/logger"
	pkg_http_middleware "couriers/pkg/transport/http/middleware"
	"net/http"

	"go.uber.org/zap"
)

func Logger(logger *pkg_logger.Logger) pkg_http_middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(pkg_http_middleware.RequestIDKey)

			log := logger.With(zap.String("request_ID", requestID), zap.String("URL", r.URL.String()))

			ctx := pkg_logger.ToContext(r.Context(), log)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
