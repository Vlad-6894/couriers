package pkg_kafka_producer

import "time"

type Producer interface {
	GetAddr() string
	GetTopic() string
	GetMaxAttempts() int
	GetWriteTimeout() time.Duration
	GetReadTimeout() time.Duration
}
