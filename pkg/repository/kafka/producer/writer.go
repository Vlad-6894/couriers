package pkg_kafka_producer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type KafkaWriter struct {
	*kafka.Writer
}

func NewKafkaWriter(config KafkaProducerConfig) *KafkaWriter {
	return &KafkaWriter{
		Writer: &kafka.Writer{
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
