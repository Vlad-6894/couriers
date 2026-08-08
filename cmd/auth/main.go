package main

import (
	"context"
	auth_repository_postgres "couriers/internal/Auth/features/repository/postgres"
	auth_service "couriers/internal/Auth/features/service"
	auth_http_transport "couriers/internal/Auth/features/transport/http"
	pkg_jwt "couriers/pkg/jwt"
	pkg_logger "couriers/pkg/logger"
	pkg_postgres_pool "couriers/pkg/repository/postgres/pool"
	pkg_http_middleware "couriers/pkg/transport/http/middleware"
	pkg_http_server "couriers/pkg/transport/http/server"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := pkg_logger.NewLogger(pkg_logger.NewLoggerConfigMust())
	if err != nil {
		fmt.Println("init logger error!")
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("init connection pool!")
	pool, err := pkg_postgres_pool.NewPostgresConnectionPool(ctx, pkg_postgres_pool.NewPostgresConnectionConfigMust())
	if err != nil {
		logger.Fatal("connection pool init error!")
	}
	defer pool.Close()

	logger.Debug("init features")
	authRepository := auth_repository_postgres.NewAuthRepository(pool)
	authService := auth_service.NewAuthService(authRepository)
	authTransportHTTP := auth_http_transport.NewAuthHTTPHandler(authService)

	logger.Debug("init http server!")

	httpServer := pkg_http_server.NewHTTPServer(
		pkg_http_server.NewHTTPServerConfigMust(),
		logger,
		pkg_http_middleware.CORS(),
		pkg_http_middleware.RequestID(),
		pkg_http_middleware.LoggerForAuthService(logger),
		pkg_http_middleware.Trace(),
		pkg_http_middleware.Panic(),
	)

	apiVersionRouterUsers := pkg_http_server.NewApiVersionRouter(pkg_http_server.ApiVersion1, pkg_jwt.User)
	apiVersionRouterUsers.RegisterRoutes(authTransportHTTP.UserRoutes()...)

	apiVersionRouterCouriers := pkg_http_server.NewApiVersionRouter(pkg_http_server.ApiVersion1, pkg_jwt.Courier)
	apiVersionRouterCouriers.RegisterRoutes(authTransportHTTP.CourierRoutes()...)

	httpServer.RegisterRouters(apiVersionRouterUsers, apiVersionRouterCouriers)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Error run HTTP server!", zap.Error(err))
	}
}
