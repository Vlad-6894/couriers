package dispetch_kafka_transport

import (
	"context"
	dispetch_domains "couriers/internal/Dispetch/core/domains"
	pkg_logger "couriers/pkg/logger"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type DispetchOrdersKafkaConsumer struct {
	*kafka.Reader
	dispetchOrdersService DispetchOrdersService
	log                   *pkg_logger.Logger
}

type DispetchOrdersService interface {
	CheckUnique(
		ctx context.Context,
		orderId int,
	) (bool, error)

	DispetchOrder(
		ctx context.Context,
		order dispetch_domains.Order,
	) error
}

var (
	networkProto = "tcp"
	waitTime     = 3 * time.Second
)

func NewDispetchOrdersKafkaConsumer(
	reader *kafka.Reader,
	dispetchOrdersService DispetchOrdersService,
	log *pkg_logger.Logger,
) *DispetchOrdersKafkaConsumer {
	return &DispetchOrdersKafkaConsumer{
		Reader:                reader,
		dispetchOrdersService: dispetchOrdersService,
		log:                   log,
	}
}

func (c *DispetchOrdersKafkaConsumer) Start(
	ctx context.Context,
	brokers []string,
	topic string,
) error {
	conn, err := kafka.Dial(networkProto, brokers[0])
	if err != nil {
		c.log.Error("fail to conn to broker", zap.Error(err))
		return fmt.Errorf("fail to conn to broker: %w", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		c.log.Error("fail to count partitions ", zap.Error(err))
		return fmt.Errorf("fail to count partitions: %w", err)
	}

	numWorkers := len(partitions)
	if numWorkers == 0 {
		c.log.Error("num partitions can not be 0! ")
		return fmt.Errorf("num partitions can not be 0: %w", err)
	}

	for i := 1; i <= numWorkers; i++ {
		go c.readPartition(ctx)
	}

	return nil
}

func (c *DispetchOrdersKafkaConsumer) readPartition(ctx context.Context) {
	for {
		msg, err := c.Reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				c.log.Warn("orders consumer is cancel!")
				return
			}

			c.log.Error("orders consumer fetch message error: ", zap.Error(err))
			continue
		}

		var event OrderCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.log.Error("fail unmarshal message", zap.Error(err))
			if err := c.CommitMessages(ctx, msg); err != nil {
				c.log.Error("fail unmarshal message", zap.Error(err))
				return
			}
			continue
		}

		isUnique, err := c.dispetchOrdersService.CheckUnique(ctx, event.ID)
		if err != nil {
			c.log.Error("fail check unique message", zap.Error(err))
			continue
		}

		if !isUnique {
			c.log.Warn("idemtece! Found a dublicate!")
			if err := c.CommitMessages(ctx, msg); err != nil {
				c.log.Error("fail commit message", zap.Error(err))
				return
			}
			continue
		}

		order := orderDomainFromEventOrderCreated(event)

		if err := c.dispetchOrdersService.DispetchOrder(ctx, order); err != nil {
			c.log.Error("error dispetch order from service: ", zap.Error(err))
			time.Sleep(waitTime)
			continue
		}

		if err := c.CommitMessages(ctx, msg); err != nil {
			c.log.Error("fail commit message", zap.Error(err))
			return
		}
	}
}

func orderDomainFromEventOrderCreated(event OrderCreatedEvent) dispetch_domains.Order {
	order := dispetch_domains.NewOrder(
		event.ID,
		event.Version,
		event.Name,
		event.Price,
		event.IsComplete,
		event.City,
		event.UserID,
	)

	return order
}
