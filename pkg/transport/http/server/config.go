package pkg_http_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type HTTPServerConfig struct {
	Addr            string        `envconfig:"ADDR" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" required:"true"`
}

func NewHTTPServerConfig() (HTTPServerConfig, error) {
	var config HTTPServerConfig

	if err := envconfig.Process("HTTP", &config); err != nil {
		return HTTPServerConfig{}, fmt.Errorf("envconfig process error: %w", err)
	}

	return config, nil
}

func NewHTTPServerConfigMust() HTTPServerConfig {
	config, err := NewHTTPServerConfig()
	if err != nil {
		panic(err)
	}

	return config
}
