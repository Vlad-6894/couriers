//go:build integration

package couriers_repository_kafka

import (
	"context"
	courier_domains "couriers/internal/Couriers/core/domains"
	dispetch_core_transport_kafka "couriers/internal/Dispetch/core/transport/kafka"
	dispetch_core_kafka_transport_topik "couriers/internal/Dispetch/core/transport/kafka/confirm/topik"
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
	topic      = "confirm_orders"
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
	os.Setenv("KAFKA_CONFIRM_TOPIC", topic)
	os.Setenv("KAFKA_NUM_PARTITIONS", "4")
	os.Setenv("KAFKA_CREATE_TIME", "10s")
	os.Setenv("KAFKA_PRODUCER_ADDRESS", kafkaAddress)
	os.Setenv("KAFKA_PRODUCER_TOPIC", topic)
	os.Setenv("KAFKA_PRODUCER_WRITE_TIMEOUT", "20s")
	os.Setenv("KAFKA_PRODUCER_MAX_ATTEMPTS", "5")
	os.Setenv("KAFKA_PRODUCER_READ_TIMEOUT", "20s")
	os.Setenv("KAFKA_CONSUMER_ADDRESS", kafkaAddress)
	os.Setenv("KAFKA_CONSUMER_ORDERS_TOPIC", "orders")
	os.Setenv("KAFKA_CONSUMER_ORDERS_GROUP_ID", "rdsscvhgg")
	os.Setenv("KAFKA_CONSUMER_COMMIT_INTERVAL", "0s")
	os.Setenv("KAFKA_CONSUMER_CONFIRM_TOPIC", topic)
	os.Setenv("KAFKA_CONSUMER_CONFIRM_GROUP_ID", "dispetch-service-confirmations-group")

	defer func() {
		os.Unsetenv("KAFKA_ADDRESS")
		os.Unsetenv("KAFKA_CONFIRM_TOPIC")
		os.Unsetenv("KAFKA_NUM_PARTITIONS")
		os.Unsetenv("KAFKA_CREATE_TIME")
		os.Unsetenv("KAFKA_PRODUCER_ADDRESS")
		os.Unsetenv("KAFKA_PRODUCER_TOPIC")
		os.Unsetenv("KAFKA_PRODUCER_WRITE_TIMEOUT")
		os.Unsetenv("KAFKA_PRODUCER_MAX_ATTEMPTS")
		os.Unsetenv("KAFKA_PRODUCER_READ_TIMEOUT")
		os.Unsetenv("KAFKA_CONSUMER_ADDRESS")
		os.Unsetenv("KAFKA_CONSUMER_ORDERS_TOPIC")
		os.Unsetenv("KAFKA_CONSUMER_ORDERS_GROUP_ID")
		os.Unsetenv("KAFKA_CONSUMER_COMMIT_INTERVAL")
		os.Unsetenv("KAFKA_CONSUMER_CONFIRM_TOPIC")
		os.Unsetenv("KAFKA_CONSUMER_CONFIRM_GROUP_ID")
	}()

	if err := pkg_kafka_topic.InitTopic(ctx, dispetch_core_kafka_transport_topik.NewConfirmKafkaTopicConfigMust()); err != nil {
		t.Fatalf("fail init kafka topic: %v", err)
	}

	writer := pkg_kafka_producer.NewKafkaWriter(pkg_kafka_producer.NewKafkaProducerConfigMust())
	defer writer.Close()

	reader := dispetch_core_transport_kafka.NewConfirmReader(dispetch_core_transport_kafka.NewKafkaConsumerConfigMust())
	defer reader.Close()

	producer := NewConfirmsKafkaProducer(writer)

	t.Run("test confirm topic", func(t *testing.T) {
		if err := producer.SendConfirm(ctx, courier_domains.NewDispetchedOrder(1, 1, 1)); err != nil {
			t.Errorf("fail send confirm: %v", err)
		}

		message, err := reader.FetchMessage(ctx)
		if err != nil {
			t.Errorf("fail get message: %v", err)
		}

		var event ConfirmEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			t.Errorf("fail to unmarshal message: %v", err)
		}

		if err := reader.CommitMessages(ctx, message); err != nil {
			t.Errorf("fail to commit message: %v", err)
		}
	})
}
