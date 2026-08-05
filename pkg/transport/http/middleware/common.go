package pkg_http_middleware

import (
	pkg_jwt "couriers/pkg/jwt"
	pkg_logger "couriers/pkg/logger"
	pkg_http_response "couriers/pkg/transport/http/response"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func Auth() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secretKey := os.Getenv(pkg_jwt.JwtKey)
			if secretKey == "" {
				panic("Fatal! Fail to get jwt secret key!")
			}

			jwtSecret := []byte(secretKey)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "There is no Token!", http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "Invalid token format!", http.StatusUnauthorized)
				return
			}
			tokenString := tokenParts[0]

			claims := &pkg_jwt.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}

				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "token does not valid!", http.StatusUnauthorized)
				return
			}

			path := r.URL.Path

			if strings.HasPrefix(path, "api/v1/user") && claims.Role != pkg_jwt.User {
				http.Error(w, "you are not user, but you try to use user's funcs!", http.StatusUnauthorized)
				return
			}

			if strings.HasPrefix(path, "api/v1/courier") && claims.Role != pkg_jwt.Courier {
				http.Error(w, "you are not courier, but you try to use courier's funcs!", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = pkg_jwt.PersonIDToContext(ctx, claims.PersonID)
			ctx = pkg_jwt.RoleToContext(ctx, claims.Role)
			ctx = pkg_jwt.CityToContext(ctx, claims.City)

			next.ServeHTTP(w, r.WithContext(ctx))
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

func LoggerForAuthService(logger *pkg_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDKey)

			log := logger.With(zap.String("request_ID", requestID), zap.String("URL", r.URL.String()))

			ctx := pkg_logger.ToContext(r.Context(), log)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Logger(logger *pkg_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			personID := pkg_jwt.PersonIDFromContext(ctx)
			role := pkg_jwt.RoleFromContext(ctx)
			city := pkg_jwt.CityFromContext(ctx)

			personIDStr := strconv.Itoa(personID)

			requestID := r.Header.Get(RequestIDKey)

			log := logger.With(
				zap.String("request_ID", requestID),
				zap.String("URL", r.URL.String()),
				zap.String("person_ID", personIDStr),
				zap.String("role", role),
				zap.String("city", city),
			)

			ctx = pkg_logger.ToContext(ctx, log)

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
