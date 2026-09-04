//go:build integration

package dispetch_kafka_repository

import (
	"context"
	dispetch_core_kafka_producer "couriers/internal/Dispetch/core/repository/kafka/producer"
	dispetch_core_kafka_repozitory_topic "couriers/internal/Dispetch/core/repository/kafka/topic"
	dispetch_core_transport_kafka "couriers/internal/Dispetch/core/transport/kafka"
	dispetch_kafka_transport "couriers/internal/Dispetch/features/order/transport/kafka"
	pkg_kafka_producer "couriers/pkg/repository/kafka/producer"
	pkg_kafka_topic "couriers/pkg/repository/kafka/topic"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	tc_kafka "github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	kafkaImage = "confluentinc/cp-kafka:7.6.0"
	topic      = "dispetched_orders"
)

func TestBroker_integration(t *testing.T) {
	ctx := t.Context()

	kafkaContainer, err := tc_kafka.Run(
		ctx,
		kafkaImage,
		tc_kafka.WithClusterID("Adrfcstggg"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("9093/tcp"),
		),
	)
	if err != nil {
		t.Errorf("fail run kafka container: %v", err)
	}

	defer func() {
		if err := kafkaContainer.Terminate(context.Background()); err != nil {
			fmt.Println("fail to close kafka container: ", err)
		}
	}()

	brokers, err := kafkaContainer.Brokers(ctx)
	if err != nil {
		t.Fatalf("fail get kafka container addresses: %v", err)
	}

	kafkaAddress := brokers[0]

	os.Setenv("KAFKA_ADDRESS", kafkaAddress)
	os.Setenv("KAFKA_ORDERS_DISPETCHED_TOPIC", topic)
	os.Setenv("KAFKA_NUM_PARTITIONS", "4")
	os.Setenv("KAFKA_CREATE_TIME", "10s")
	os.Setenv("KAFKA_WRITE_TIMEOUT", "20s")
	os.Setenv("KAFKA_MAX_ATTEMPTS", "5")
	os.Setenv("KAFKA_READ_TIMEOUT", "20s")
	os.Setenv("KAFKA_CONSUMER_ADDRESS", kafkaAddress)
	os.Setenv("KAFKA_CONSUMER_ORDERS_TOPIC", topic)
	os.Setenv("KAFKA_CONSUMER_ORDERS_GROUP_ID", "test-service-dispetched-group")
	os.Setenv("KAFKA_CONSUMER_COMMIT_INTERVAL", "0s")
	os.Setenv("KAFKA_CONSUMER_CONFIRM_TOPIC", "confirm")
	os.Setenv("KAFKA_CONSUMER_CONFIRM_GROUP_ID", "dgghffvvd")

	defer func() {
		os.Unsetenv("KAFKA_ADDRESS")
		os.Unsetenv("KAFKA_ORDERS_DISPETCHED_TOPIC")
		os.Unsetenv("KAFKA_NUM_PARTITIONS")
		os.Unsetenv("KAFKA_CREATE_TIME")
		os.Unsetenv("KAFKA_WRITE_TIMEOUT")
		os.Unsetenv("KAFKA_MAX_ATTEMPTS")
		os.Unsetenv("KAFKA_READ_TIMEOUT")
		os.Unsetenv("KAFKA_CONSUMER_ADDRESS")
		os.Unsetenv("KAFKA_CONSUMER_ORDERS_TOPIC")
		os.Unsetenv("KAFKA_CONSUMER_ORDERS_GROUP_ID")
		os.Unsetenv("KAFKA_CONSUMER_COMMIT_INTERVAL")
		os.Unsetenv("KAFKA_CONSUMER_CONFIRM_TOPIC")
		os.Unsetenv("KAFKA_CONSUMER_CONFIRM_GROUP_ID")
	}()

	if err := pkg_kafka_topic.InitTopic(ctx, dispetch_core_kafka_repozitory_topic.NewDispetchedOrdersKafkaTopicConfigMust()); err != nil {
		t.Fatalf("fail init kafka topic: %v", err)
	}

	writer := pkg_kafka_producer.NewKafkaWriter(dispetch_core_kafka_producer.NewDispetchedKafkaProducerConfigMust())
	defer writer.Close()

	reader := dispetch_core_transport_kafka.NewOrderReader(dispetch_core_transport_kafka.NewKafkaConsumerConfigMust())
	defer reader.Close()

	producer := NewDispetchKafkaRepository(writer)

	t.Run("send order test", func(t *testing.T) {
		if err := producer.SendDispetchedOrder(ctx, 1, 1, 1, "Moscow"); err != nil {
			t.Errorf("fail send order: %v", err)
		}

		message, err := reader.FetchMessage(ctx)
		if err != nil {
			t.Errorf("fail get message: %v", err)
		}

		var event dispetch_kafka_transport.OrderCreatedEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			t.Errorf("fail to unmarshal message: %v", err)
		}

		if err := reader.CommitMessages(ctx, message); err != nil {
			t.Errorf("fail to commit message: %v", err)
		}
	})
}
