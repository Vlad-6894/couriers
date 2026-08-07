package pkg_kafka_producer

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type KafkaProducerConfig struct {
	Addr         string        `envconfig:"ADDRESS" required:"true"`
	Topic        string        `envconfig:"TOPIC" required:"true"`
	MaxAttempts  int           `envconfig:"MAX_ATTEMPTS" required:"true"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" required:"true"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" required:"true"`
}

func NewKafkaProducerConfig() (KafkaProducerConfig, error) {
	var config KafkaProducerConfig

	if err := envconfig.Process("KAFKA_PRODUCER", &config); err != nil {
		return KafkaProducerConfig{}, fmt.Errorf("fail to proccess kafka orders config: %w", err)
	}

	return config, nil
}

func NewKafkaProducerConfigMust() KafkaProducerConfig {
	config, err := NewKafkaProducerConfig()
	if err != nil {
		panic(err)
	}

	return config
}
