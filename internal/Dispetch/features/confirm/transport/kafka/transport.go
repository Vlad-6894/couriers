package dispetch_kafka_confirm_transport

import (
	"context"
	dispetch_domains "couriers/internal/Dispetch/core/domains"
	pkg_logger "couriers/pkg/logger"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type ConfirmKafkaConsumer struct {
	*kafka.Reader
	dispetchConfirmService DispetchConfirmService
	log                    *pkg_logger.Logger
	wg                     *sync.WaitGroup
}

type DispetchConfirmService interface {
	ConfirmOrder(
		ctx context.Context,
		confirm dispetch_domains.Confirm,
	) error
}

var (
	networkProto = "tcp"
	waitTime     = 3 * time.Second
)

func NewConfirmKafkaConsumer(
	reader *kafka.Reader,
	dispetchConfirmService DispetchConfirmService,
	log *pkg_logger.Logger,
	wg *sync.WaitGroup,
) *ConfirmKafkaConsumer {
	return &ConfirmKafkaConsumer{
		Reader:                 reader,
		dispetchConfirmService: dispetchConfirmService,
		log:                    log,
		wg:                     wg,
	}
}

func (c *ConfirmKafkaConsumer) Start(
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

func (c *ConfirmKafkaConsumer) readPartition(ctx context.Context) {
	for {
		msg, err := c.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				c.log.Warn("confirm consumer is cancel!")
				return
			}

			c.log.Error("confirm consumer fetch message error: ", zap.Error(err))
			continue
		}

		var event ConfirmEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.log.Error("fail unmarshal message", zap.Error(err))
			if err := c.CommitMessages(ctx, msg); err != nil {
				c.log.Error("fail unmarshal message", zap.Error(err))
				return
			}
			continue
		}

		confirm := confirmDomainFromEvent(event)

		if err := c.dispetchConfirmService.ConfirmOrder(ctx, confirm); err != nil {
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

func confirmDomainFromEvent(event ConfirmEvent) dispetch_domains.Confirm {
	confirm := dispetch_domains.NewConfirm(event.OrderID, event.CourierID)
	return confirm
}
