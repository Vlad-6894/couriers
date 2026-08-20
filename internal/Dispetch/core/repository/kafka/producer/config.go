package dispetch_core_kafka_producer

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type DispetchKafkaProducerConfig struct {
	Addr         string        `envconfig:"ADDRESS" required:"true"`
	Topic        string        `envconfig:"ORDERS_DISPETCHED_TOPIC" required:"true"`
	MaxAttempts  int           `envconfig:"MAX_ATTEMPTS" required:"true"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" required:"true"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" required:"true"`
}

func NewDispetchedKafkaProducerConfig() (DispetchKafkaProducerConfig, error) {
	var config DispetchKafkaProducerConfig

	if err := envconfig.Process("KAFKA", &config); err != nil {
		return DispetchKafkaProducerConfig{}, fmt.Errorf("fail to proccess kafka dispetch config: %w", err)
	}

	return config, nil
}

func NewDispetchedKafkaProducerConfigMust() DispetchKafkaProducerConfig {
	config, err := NewDispetchedKafkaProducerConfig()
	if err != nil {
		panic(err)
	}

	return config
}

func (p DispetchKafkaProducerConfig) GetAddr() string {
	return p.Addr
}

func (p DispetchKafkaProducerConfig) GetTopic() string {
	return p.Topic
}

func (p DispetchKafkaProducerConfig) GetMaxAttempts() int {
	return p.MaxAttempts
}

func (p DispetchKafkaProducerConfig) GetWriteTimeout() time.Duration {
	return p.WriteTimeout
}

func (p DispetchKafkaProducerConfig) GetReadTimeout() time.Duration {
	return p.ReadTimeout
}
