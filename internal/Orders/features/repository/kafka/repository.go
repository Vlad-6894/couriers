package orders_repository_kafka

import (
	orders_domains "couriers/internal/Orders/core/domains"
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

type OrderEvent struct {
	ID         int    `json:"order_id"`
	Version    int    `json:"order_version"`
	Name       string `json:"order_name"`
	Price      int    `json:"order_price"`
	IsComplete bool   `json:"is_complete"`
	City       string `json:"city"`
	UserID     int    `json:"user_id"`
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

func newOrderEventFromDomain(order orders_domains.Order) OrderEvent {
	return OrderEvent{
		ID:         order.ID,
		Version:    order.Version,
		Name:       order.Name,
		Price:      order.Price,
		IsComplete: order.IsComplete,
		City:       order.City,
		UserID:     order.UserID,
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
