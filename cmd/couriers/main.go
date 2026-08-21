package main

import (
	"context"
	couriers_core_kafka_transport "couriers/internal/Couriers/core/transport/kafka"
	couriers_repository_kafka "couriers/internal/Couriers/features/repository/kafka"
	couriers_redis_repository "couriers/internal/Couriers/features/repository/redis"
	couriers_service "couriers/internal/Couriers/features/service"
	couriers_http_transport "couriers/internal/Couriers/features/transport/http"
	couriers_kafka_transport "couriers/internal/Couriers/features/transport/kafka"
	pkg_jwt "couriers/pkg/jwt"
	pkg_logger "couriers/pkg/logger"
	pkg_kafka_producer "couriers/pkg/repository/kafka/producer"
	pkg_repository_redis "couriers/pkg/repository/redis"
	pkg_http_middleware "couriers/pkg/transport/http/middleware"
	pkg_http_server "couriers/pkg/transport/http/server"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var (
	errorsChanSize = 1
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	wg := &sync.WaitGroup{}

	logger, err := pkg_logger.NewLogger(pkg_logger.NewLoggerConfigMust())
	if err != nil {
		fmt.Println("init logger error!")
		os.Exit(1)
	}
	defer logger.Close()

	clientRedis, err := pkg_repository_redis.NewRedisClient(pkg_repository_redis.NewRedisConfigMust())
	if err != nil {
		logger.Fatal("fail init redis clent: ", zap.Error(err))
	}

	writer := pkg_kafka_producer.NewKafkaWriter(
		pkg_kafka_producer.NewKafkaProducerConfigMust(),
	)
	defer writer.Close()

	consumerConfig := couriers_core_kafka_transport.NewKafkaConsumerConfigMust()

	orderReader := couriers_core_kafka_transport.NewOrderReader(
		consumerConfig,
	)
	defer orderReader.Close()

	logger.Debug("Init feature")

	cache := couriers_redis_repository.NewCouriersCache(clientRedis)
	confirmProducer := couriers_repository_kafka.NewConfirmsKafkaProducer(writer)

	service := couriers_service.NewCouriersService(cache, confirmProducer)

	ordersConsumer := couriers_kafka_transport.NewCouriersKafkaConsumer(
		orderReader,
		service,
		logger,
		wg,
	)

	courierTransport := couriers_http_transport.NewCouriersHTTPHandler(service)

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

	apiVersionRouterCouriers := pkg_http_server.NewApiVersionRouter(pkg_http_server.ApiVersion1, pkg_jwt.Courier)
	apiVersionRouterCouriers.RegisterRoutes(courierTransport.Routes()...)

	httpServer.RegisterRouters(apiVersionRouterCouriers)

	if err := ordersConsumer.Start(
		ctx,
		consumerConfig.Brokers,
		consumerConfig.OrdersTopic,
	); err != nil {
		logger.Fatal("fail to start orders consumer: ", zap.Error(err))
	}

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Error run HTTP server!", zap.Error(err))
		cancel()
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("program stopped")
	case <-time.After(5 * time.Second):
		logger.Error("time is up!")
	}
}
