package couriers_repository_kafka

import (
	"context"
	courier_domains "couriers/internal/Couriers/core/domains"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func (p *ConfirmsKafkaProducer) SendConfirm(
	ctx context.Context,
	orderInfo courier_domains.DispetchedOrder,
) error {
	event := eventFromDomain(orderInfo)

	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error message encode to byte: %w", err)
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Value: message,
	}); err != nil {
		return fmt.Errorf("fail to write message: %w", err)
	}

	return nil
}

func eventFromDomain(order courier_domains.DispetchedOrder) ConfirmEvent {
	event := ConfirmEvent{
		OrderID:   order.OrderID,
		CourierID: order.CourierID,
	}

	return event
}
