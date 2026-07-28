package pkg_postgres_pool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type PostgresConnectionPoolConfig struct {
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" default:"5432"`
	User     string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Database string        `envconfig:"DB" requred:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" requred:"true"`
}

func NewPostgresConnectionConfig() (PostgresConnectionPoolConfig, error) {
	var config PostgresConnectionPoolConfig

	if err := envconfig.Process("POSTGRES", &config); err != nil {
		return PostgresConnectionPoolConfig{}, fmt.Errorf("Proccess config error: %w", err)
	}

	return config, nil
}

func NewPostgresConnectionConfigMust() PostgresConnectionPoolConfig {
	config, err := NewPostgresConnectionConfig()
	if err != nil {
		panic(err)
	}

	return config
}
