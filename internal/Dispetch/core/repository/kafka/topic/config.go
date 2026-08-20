package dispetch_core_kafka_repozitory_topic

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type DispetchedOrdersKafkaTopicConfig struct {
	BrokerAddrr   string        `envconfig:"ADDRESS" required:"true"`
	TopicName     string        `envconfig:"ORDERS_DISPETCHED_TOPIC" required:"true"`
	NumPartitions int           `envconfig:"NUM_PARTITIONS" required:"true"`
	CreateTime    time.Duration `envconfig:"CREATE_TIME" required:"true"`
}

func NewDispetchedOrdersKafkaTopicConfig() (DispetchedOrdersKafkaTopicConfig, error) {
	var config DispetchedOrdersKafkaTopicConfig

	if err := envconfig.Process("KAFKA", &config); err != nil {
		return DispetchedOrdersKafkaTopicConfig{}, fmt.Errorf("fail to proccess kafka topic config: %w", err)
	}

	return config, nil
}

func NewDispetchedOrdersKafkaTopicConfigMust() DispetchedOrdersKafkaTopicConfig {
	config, err := NewDispetchedOrdersKafkaTopicConfig()
	if err != nil {
		panic(err)
	}

	return config
}

func (c DispetchedOrdersKafkaTopicConfig) GetBrokerAddrr() string {
	return c.BrokerAddrr
}

func (c DispetchedOrdersKafkaTopicConfig) GetTopicName() string {
	return c.TopicName
}

func (c DispetchedOrdersKafkaTopicConfig) GetNumPartitions() int {
	return c.NumPartitions
}

func (c DispetchedOrdersKafkaTopicConfig) GetCreateTime() time.Duration {
	return c.CreateTime
}
