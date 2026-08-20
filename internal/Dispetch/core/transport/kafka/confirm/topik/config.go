package dispetch_core_kafka_transport_topik

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type ConfirmKafkaTopicConfig struct {
	BrokerAddrr   string        `envconfig:"ADDRESS" required:"true"`
	TopicName     string        `envconfig:"CONFIRM_TOPIC" required:"true"`
	NumPartitions int           `envconfig:"NUM_PARTITIONS" required:"true"`
	CreateTime    time.Duration `envconfig:"CREATE_TIME" required:"true"`
}

func NewConfirmKafkaTopicConfig() (ConfirmKafkaTopicConfig, error) {
	var config ConfirmKafkaTopicConfig

	if err := envconfig.Process("KAFKA", &config); err != nil {
		return ConfirmKafkaTopicConfig{}, fmt.Errorf("fail to proccess kafka topic config: %w", err)
	}

	return config, nil
}

func NewConfirmKafkaTopicConfigMust() ConfirmKafkaTopicConfig {
	config, err := NewConfirmKafkaTopicConfig()
	if err != nil {
		panic(err)
	}

	return config
}

func (c ConfirmKafkaTopicConfig) GetBrokerAddrr() string {
	return c.BrokerAddrr
}

func (c ConfirmKafkaTopicConfig) GetTopicName() string {
	return c.TopicName
}

func (c ConfirmKafkaTopicConfig) GetNumPartitions() int {
	return c.NumPartitions
}

func (c ConfirmKafkaTopicConfig) GetCreateTime() time.Duration {
	return c.CreateTime
}
