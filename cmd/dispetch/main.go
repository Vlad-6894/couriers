package main

import (
	"context"
	dispetch_core_kafka_producer "couriers/internal/Dispetch/core/repository/kafka/producer"
	dispetch_core_kafka_repozitory_topic "couriers/internal/Dispetch/core/repository/kafka/topic"
	dispetch_core_transport_kafka "couriers/internal/Dispetch/core/transport/kafka"
	dispetch_core_kafka_transport_topik "couriers/internal/Dispetch/core/transport/kafka/confirm/topik"
	dispetch_confirm_repository_postgres "couriers/internal/Dispetch/features/confirm/repozitory/postgres"
	dispetch_confirm_service "couriers/internal/Dispetch/features/confirm/service"
	dispetch_kafka_confirm_transport "couriers/internal/Dispetch/features/confirm/transport/kafka"
	dispetch_kafka_repository "couriers/internal/Dispetch/features/order/repository/kafka"
	dispetch_postgres_repository "couriers/internal/Dispetch/features/order/repository/postgres"
	dispetch_redis_repository "couriers/internal/Dispetch/features/order/repository/redis"
	dispetch_service "couriers/internal/Dispetch/features/order/service"
	dispetch_kafka_transport "couriers/internal/Dispetch/features/order/transport/kafka"
	pkg_logger "couriers/pkg/logger"
	pkg_kafka_producer "couriers/pkg/repository/kafka/producer"
	pkg_kafka_topic "couriers/pkg/repository/kafka/topic"
	pkg_postgres_pool "couriers/pkg/repository/postgres/pool"
	pkg_repository_redis "couriers/pkg/repository/redis"
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
	errorsChan := make(chan error, errorsChanSize)

	logger, err := pkg_logger.NewLogger(pkg_logger.NewLoggerConfigMust())
	if err != nil {
		fmt.Println("init logger error!")
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("init connection pool!")
	pool, err := pkg_postgres_pool.NewPostgresConnectionPool(ctx, pkg_postgres_pool.NewPostgresConnectionConfigMust())
	if err != nil {
		logger.Fatal("connection pool init error: ", zap.Error(err))
	}
	defer pool.Close()

	clientRedis, err := pkg_repository_redis.NewRedisClient(pkg_repository_redis.NewRedisConfigMust())
	if err != nil {
		logger.Fatal("fail init redis clent: ", zap.Error(err))
	}

	if err := pkg_kafka_topic.InitTopic(
		ctx,
		dispetch_core_kafka_repozitory_topic.NewDispetchedOrdersKafkaTopicConfigMust(),
	); err != nil {
		logger.Error("fail to init topic:", zap.Error(err))
	}

	if err := pkg_kafka_topic.InitTopic(
		ctx,
		dispetch_core_kafka_transport_topik.NewConfirmKafkaTopicConfigMust(),
	); err != nil {
		logger.Error("fail to init topic:", zap.Error(err))
	}

	writer := pkg_kafka_producer.NewKafkaWriter(
		dispetch_core_kafka_producer.NewDispetchedKafkaProducerConfigMust(),
	)
	defer writer.Close()

	kafkaConsumerConfig := dispetch_core_transport_kafka.NewKafkaConsumerConfigMust()

	ordersReader := dispetch_core_transport_kafka.NewOrderReader(
		kafkaConsumerConfig,
	)
	defer ordersReader.Close()

	confirmReader := dispetch_core_transport_kafka.NewConfirmReader(
		kafkaConsumerConfig,
	)

	logger.Debug("Init feature orders")

	ordersPostgres := dispetch_postgres_repository.NewDispetchRepositoryPostgres(pool)
	ordersCache := dispetch_redis_repository.NewRedisRepository(clientRedis)
	ordersProducer := dispetch_kafka_repository.NewDispetchKafkaRepository(writer)

	ordersService := dispetch_service.NewOrdersDispetchService(
		ordersPostgres,
		ordersCache,
		ordersProducer,
	)

	ordersConsumer := dispetch_kafka_transport.NewDispetchOrdersKafkaConsumer(
		ordersReader,
		ordersService,
		logger,
		wg,
	)

	logger.Debug("Init feature confirm!")

	confirmPostgres := dispetch_confirm_repository_postgres.NewConfirmPostgresDatabase(pool)

	confirmService := dispetch_confirm_service.NewDispetchConfirmService(confirmPostgres)

	confirmConsumer := dispetch_kafka_confirm_transport.NewConfirmKafkaConsumer(
		confirmReader,
		confirmService,
		logger,
		wg,
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := ordersService.UpdateCache(ctx); err != nil {
			errorsChan <- err
		}
	}()

	if err := ordersConsumer.Start(
		ctx,
		kafkaConsumerConfig.Brokers,
		kafkaConsumerConfig.OrdersTopic,
	); err != nil {
		logger.Fatal("fail to start orders consumer: ", zap.Error(err))
	}

	if err := confirmConsumer.Start(
		ctx,
		kafkaConsumerConfig.Brokers,
		kafkaConsumerConfig.ConfirmTopic,
	); err != nil {
		logger.Fatal("fail to start confirm consumer: ", zap.Error(err))
	}

	select {
	case <-ctx.Done():
		logger.Info("stop program by syscall")
	case err := <-errorsChan:
		logger.Error("stop program by error:", zap.Error(err))
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
