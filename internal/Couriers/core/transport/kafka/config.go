package couriers_core_kafka_transport

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type KafkaConsumerConfig struct {
	Brokers        []string      `envconfig:"ADDRESS" required:"true"`
	OrdersTopic    string        `envconfig:"ORDERS_TOPIC" required:"true"`
	OrdersGroupID  string        `envconfig:"ORDERS_GROUP_ID" required:"true"`
	CommitInterval time.Duration `envconfig:"COMMIT_INTERVAL" default:"0s" `
}

func NewKafkaConsumerConfig() (KafkaConsumerConfig, error) {
	var config KafkaConsumerConfig

	if err := envconfig.Process("KAFKA_CONSUMER", &config); err != nil {
		return KafkaConsumerConfig{}, fmt.Errorf("fail to process kafka config: %w", err)
	}

	return config, nil
}

func NewKafkaConsumerConfigMust() KafkaConsumerConfig {
	config, err := NewKafkaConsumerConfig()
	if err != nil {
		panic(err)
	}

	return config
}
