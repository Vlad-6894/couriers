package couriers_kafka_transport

import (
	"context"
	courier_domains "couriers/internal/Couriers/core/domains"
	couriers_transport "couriers/internal/Couriers/features/transport"
	pkg_logger "couriers/pkg/logger"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type CouriersKafkaConsumer struct {
	*kafka.Reader
	service couriers_transport.CouriersService
	log     *pkg_logger.Logger
	wg      *sync.WaitGroup
}

var (
	networkProto = "tcp"
	waitTime     = 3 * time.Second
)

func NewCouriersKafkaConsumer(
	reader *kafka.Reader,
	service couriers_transport.CouriersService,
	log *pkg_logger.Logger,
	wg *sync.WaitGroup,
) *CouriersKafkaConsumer {
	return &CouriersKafkaConsumer{
		Reader:  reader,
		service: service,
		log:     log,
		wg:      wg,
	}
}

func (c *CouriersKafkaConsumer) Start(
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
		c.wg.Add(1)
		go c.readPartition(ctx)
	}

	return nil
}

func (c *CouriersKafkaConsumer) readPartition(ctx context.Context) {
	defer c.wg.Done()
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

		var event DispetchedOrderEventForCourier
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.log.Error("fail unmarshal message", zap.Error(err))
			if err := c.CommitMessages(ctx, msg); err != nil {
				c.log.Error("fail unmarshal message", zap.Error(err))
				return
			}
			continue
		}

		order := dispetchedOrderDomainFromEvent(event)

		if err := c.service.SaveToCache(ctx, order); err != nil {
			c.log.Error("error save order from service: ", zap.Error(err))
			time.Sleep(waitTime)
			continue
		}

		if err := c.CommitMessages(ctx, msg); err != nil {
			c.log.Error("fail commit message", zap.Error(err))
			return
		}
	}
}

func dispetchedOrderDomainFromEvent(
	event DispetchedOrderEventForCourier,
) courier_domains.DispetchedOrder {
	order := courier_domains.NewDispetchedOrder(
		event.OrderID,
		event.Version,
		event.CourierID,
	)

	return order
}
