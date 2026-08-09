package pkg_repository_redis

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type RedisConfig struct {
	RedisHost     string        `envconfig:"HOST" required:"true"`
	RedisPort     string        `envconfig:"PORT" required:"true"`
	RedisPassword string        `envconfig:"PASSWORD" required:"true"`
	RedisTimeout  time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func NewRedisConfig() (RedisConfig, error) {
	var config RedisConfig

	if err := envconfig.Process("REDIS", &config); err != nil {
		return RedisConfig{}, fmt.Errorf("fail to process redis config: %w", err)
	}

	return config, nil
}

func NewRedisConfigMust() RedisConfig {
	config, err := NewRedisConfig()
	if err != nil {
		panic(err)
	}

	return config
}
