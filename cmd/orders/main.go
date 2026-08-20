package main

import (
	"context"
	orders_repository_kafka "couriers/internal/Orders/features/repository/kafka"
	orders_repository_posgres "couriers/internal/Orders/features/repository/postgres"
	orders_service "couriers/internal/Orders/features/service"
	orders_http_transport "couriers/internal/Orders/features/transport/http"
	pkg_jwt "couriers/pkg/jwt"
	pkg_logger "couriers/pkg/logger"
	pkg_kafka_producer "couriers/pkg/repository/kafka/producer"
	pkg_kafka_topic "couriers/pkg/repository/kafka/topic"
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

	if err := pkg_kafka_topic.InitTopic(ctx, pkg_kafka_topic.NewKafkaTopicConfigMust()); err != nil {
		logger.Error("fail to init topic:", zap.Error(err))
	}

	logger.Debug("init features!")

	writer := pkg_kafka_producer.NewKafkaWriter(
		pkg_kafka_producer.NewKafkaProducerConfigMust(),
	)
	defer writer.Close()

	producer := orders_repository_kafka.NewOrdersKafkaProducer(
		writer,
		logger,
	)

	ordersDatabase := orders_repository_posgres.NewDatabasePostgres(pool)

	ordersService := orders_service.NewOrdersService(ordersDatabase, producer)

	ordersTransport := orders_http_transport.NewOrdersHTTPHandler(ordersService)

	logger.Debug("init http server!")

	httpServer := pkg_http_server.NewHTTPServer(
		pkg_http_server.NewHTTPServerConfigMust(),
		logger,
		pkg_http_middleware.CORS(),
		pkg_http_middleware.Auth(),
		pkg_http_middleware.RequestID(),
		pkg_http_middleware.Logger(logger),
		pkg_http_middleware.Trace(),
		pkg_http_middleware.Panic(),
	)

	apiVersionRouterOrders := pkg_http_server.NewApiVersionRouter(pkg_http_server.ApiVersion1, pkg_jwt.User)
	apiVersionRouterOrders.RegisterRoutes(ordersTransport.Routes()...)

	httpServer.RegisterRouters(apiVersionRouterOrders)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Error run HTTP server!", zap.Error(err))
	}
}
