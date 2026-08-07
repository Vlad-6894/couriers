package pkg_kafka_topic

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

var (
	kafkaNetworkProto = "tcp"
)

func InitTopic(ctx context.Context, config KafkaTopicConfig) error {
	conn, err := kafka.Dial(kafkaNetworkProto, config.BrokerAddrr)
	if err != nil {
		return fmt.Errorf("fail connect to kafka cluster: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("fail get to kafka controller: %w", err)
	}

	controllerAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))

	controllerConn, err := kafka.Dial(kafkaNetworkProto, controllerAddr)
	if err != nil {
		return fmt.Errorf("fail connect to kafka controller: %w", err)
	}
	defer controllerConn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             config.TopicName,
		NumPartitions:     config.NumPartitions,
		ReplicationFactor: 1,
	}

	ctx, cancel := context.WithTimeout(ctx, config.CreateTime)
	defer cancel()

	if err := controllerConn.CreateTopics(topicConfig); err != nil {
		if errors.Is(err, kafka.TopicAlreadyExists) {
			return nil
		}
		return fmt.Errorf("fail create kafka topic: %w", err)
	}

	return nil
}
