package pkg_http_middleware

import (
	pkg_logger "couriers/pkg/logger"
	pkg_http_response "couriers/pkg/transport/http/response"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	RequestIDKey  = "X-Request-ID"
	originKey     = "Origin"
	AccessOrigin  = "Access-Control-Allow-Origin"
	AccessMethods = "Access-Control-Allow-Methods"
	AccessHeaders = "Access-Control-Allow-Headers"
)

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := map[string]struct{}{
				"http:localhost:5051": {},
			}

			origin := r.Header.Get(originKey)

			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set(AccessOrigin, origin)
				w.Header().Set(AccessMethods, "POST, GET, DELETE, PATCH, OPTIONS")
				w.Header().Set(AccessHeaders, "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDKey)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(RequestIDKey, requestID)
			w.Header().Set(RequestIDKey, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(logger *pkg_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDKey)

			log := logger.With(zap.String("request_ID", requestID), zap.String("URL", r.URL.String()))

			ctx := pkg_logger.ToContext(r.Context(), log)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := pkg_logger.FromContext(ctx)
			rw := pkg_http_response.NewResponseWriter(w)

			beforeTime := time.Now()
			log.Debug(
				"incoming HTTP Request",
				zap.String("http_method: ", r.Method),
				zap.Time("time: ", beforeTime.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"finish HTTP request!",
				zap.Int("status code: ", rw.GetStatusCode()),
				zap.Duration("latency", time.Since(beforeTime)),
			)
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := pkg_logger.FromContext(ctx)
			responseHandler := pkg_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse("Panic in the program!", p)
				}
			}()
		})
	}
}
