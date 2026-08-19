package dispetch_kafka_repository

import pkg_kafka_producer "couriers/pkg/repository/kafka/producer"

type DispetchKafkaRepository struct {
	writer pkg_kafka_producer.Writer
}

func NewDispetchKafkaRepository(
	writer pkg_kafka_producer.Writer,
) *DispetchKafkaRepository {
	return &DispetchKafkaRepository{
		writer: writer,
	}
}
