package orders_repository_kafka

import (
	"context"
	orders_domains "couriers/internal/Orders/core/domains"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func (p *OrdersKafkaProducer) SendOrder(
	ctx context.Context,
	order orders_domains.Order,
) error {
	message, err := json.Marshal(order)
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
