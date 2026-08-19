package pkg_repository_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type ClientCacheRedis interface {
	Pipeline() redis.Pipeliner
	SPop(ctx context.Context, key string) *redis.StringCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	GetTimeout() time.Duration
	Close()
}

type RedisClient struct {
	*redis.Client
	timeout time.Duration
}

func NewRedisClient(config RedisConfig) (*RedisClient, error) {
	redisAddr := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     config.RedisPassword,
		DB:           0,
		DialTimeout:  config.RedisTimeout,
		WriteTimeout: config.RedisTimeout,
		ReadTimeout:  config.RedisTimeout,
	})

	pingCtx, cancel := context.WithTimeout(context.Background(), config.RedisTimeout)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("fail to ping redis: %w", err)
	}

	client := &RedisClient{
		Client:  rdb,
		timeout: config.RedisTimeout,
	}

	return client, nil
}

func (c *RedisClient) GetTimeout() time.Duration {
	return c.timeout
}

func (c *RedisClient) Close() {
	if err := c.Client.Close(); err != nil {
		fmt.Println("fail to close redis: ", err)
	}
}
