package orders_repository_kafka

import (
	pkg_kafka_producer "couriers/pkg/repository/kafka/producer"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type OrdersKafkaProducer struct {
	writer *kafka.Writer
}

func NewOrdersKafkaProducer(config pkg_kafka_producer.KafkaProducerConfig) *OrdersKafkaProducer {
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
	}
}

func (p *OrdersKafkaProducer) Close() error {
	if err := p.writer.Close(); err != nil {
		fmt.Println("fail to close orders kafka producer!")
		return fmt.Errorf("fail to close orders kafka producer: %w", err)
	}

	return nil
}
