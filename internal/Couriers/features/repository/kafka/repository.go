package couriers_repository_kafka

import pkg_kafka_producer "couriers/pkg/repository/kafka/producer"

type ConfirmsKafkaProducer struct {
	writer pkg_kafka_producer.Writer
}

func NewConfirmsKafkaProducer(
	writer pkg_kafka_producer.Writer,
) *ConfirmsKafkaProducer {
	return &ConfirmsKafkaProducer{
		writer: writer,
	}
}
