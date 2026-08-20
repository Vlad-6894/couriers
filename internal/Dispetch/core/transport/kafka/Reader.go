package dispetch_core_transport_kafka

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

func NewConfirmReader(config KafkaConsumerConfig) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		Topic:          config.ConfirmTopic,
		GroupID:        config.ConfirmGroupID,
		CommitInterval: config.CommitInterval,
	})

	return reader
}
