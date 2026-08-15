package orders_repository_kafka

import (
	orders_domains "couriers/internal/Orders/core/domains"
	pkg_logger "couriers/pkg/logger"
	pkg_kafka_producer "couriers/pkg/repository/kafka/producer"
)

type OrdersKafkaProducer struct {
	writer pkg_kafka_producer.Writer
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

func NewOrdersKafkaProducer(writer pkg_kafka_producer.Writer, log *pkg_logger.Logger) *OrdersKafkaProducer {
	return &OrdersKafkaProducer{
		writer: writer,
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
