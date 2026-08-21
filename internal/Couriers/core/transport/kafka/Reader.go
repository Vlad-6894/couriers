package couriers_core_kafka_transport

import "github.com/segmentio/kafka-go"

func NewOrderReader(config KafkaConsumerConfig) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		Topic:          config.OrdersTopic,
		GroupID:        config.OrdersGroupID,
		CommitInterval: config.CommitInterval,
	})

	return reader
}
