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

func NewKafkaWriter(config Producer) *KafkaWriter {
	return &KafkaWriter{
		Writer: &kafka.Writer{
			Addr:         kafka.TCP(config.GetAddr()),
			Topic:        config.GetTopic(),
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
			MaxAttempts:  config.GetMaxAttempts(),
			ReadTimeout:  config.GetReadTimeout(),
			WriteTimeout: config.GetWriteTimeout(),
			Async:        false,
		},
	}
}
