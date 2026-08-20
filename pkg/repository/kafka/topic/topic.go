package pkg_kafka_topic

import "time"

type Topic interface {
	GetBrokerAddrr() string
	GetTopicName() string
	GetNumPartitions() int
	GetCreateTime() time.Duration
}
