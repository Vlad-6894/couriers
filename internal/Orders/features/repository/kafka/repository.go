package orders_repository_kafka

import (
	pkg_logger "couriers/pkg/logger"
	pkg_kafka_producer "couriers/pkg/repository/kafka/producer"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type OrdersKafkaProducer struct {
	writer *kafka.Writer
	logger *pkg_logger.Logger
}

func NewOrdersKafkaProducer(config pkg_kafka_producer.KafkaProducerConfig, log *pkg_logger.Logger) *OrdersKafkaProducer {
	return &OrdersKafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(config.Addr),
			Topic:        config.Topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
			MaxAttempts:  config.MaxAttempts,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			Async:        false,
		},
		logger: log,
	}
}

func (p *OrdersKafkaProducer) Close() error {
	if err := p.writer.Close(); err != nil {
		p.logger.Error("fail close kafka producer: %w", zap.Error(err))
		fmt.Println("fail to close orders kafka producer!")
		return fmt.Errorf("fail to close orders kafka producer: %w", err)
	}

	return nil
}
