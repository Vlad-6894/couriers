package dispetch_kafka_repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func (r *DispetchKafkaRepository) SendDispetchedOrder(
	ctx context.Context,
	orderID int,
	version int,
	courierID int,
	city string,
) error {
	event := newEvent(
		orderID,
		version,
		courierID,
	)

	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error message encode to byte: %w", err)
	}

	if err := r.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(city),
		Value: message,
	}); err != nil {
		return fmt.Errorf("fail to write message: %w", err)
	}

	return nil
}

func newEvent(
	orderID int,
	version int,
	courierID int,
) DispetchedOrderEvent {
	event := DispetchedOrderEvent{
		OrderID:   orderID,
		Version:   version,
		CourierID: courierID,
	}

	return event
}
