package pkg_kafka_topic

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type KafkaTopicConfig struct {
	BrokerAddrr   string        `envconfig:"ADDRESS" required:"true"`
	TopicName     string        `envconfig:"TOPIC" required:"true"`
	NumPartitions int           `envconfig:"NUM_PARTITIONS" required:"true"`
	CreateTime    time.Duration `envconfig:"CREATE_TIME" required:"true"`
}

func NewKafkaTopicConfig() (KafkaTopicConfig, error) {
	var config KafkaTopicConfig

	if err := envconfig.Process("KAFKA", &config); err != nil {
		return KafkaTopicConfig{}, fmt.Errorf("fail to proccess kafka topic config: %w", err)
	}

	return config, nil
}

func NewKafkaTopicConfigMust() KafkaTopicConfig {
	config, err := NewKafkaTopicConfig()
	if err != nil {
		panic(err)
	}

	return config
}
