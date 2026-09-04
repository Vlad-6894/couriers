//go:build integration

package couriers_redis_repository

import (
	"context"
	courier_domains "couriers/internal/Couriers/core/domains"
	pkg_repository_redis "couriers/pkg/repository/redis"
	"fmt"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	tc_redis "github.com/testcontainers/testcontainers-go/modules/redis"
)

var (
	redisImage        = "redis:8.6.1"
	testRedisPassword = "123"
	localHost         = "127.0.0.1"
)

func TestCache_integration(t *testing.T) {
	ctx := t.Context()

	redisContainer, err := tc_redis.Run(
		ctx,
		redisImage,
		testcontainers.WithCmd("redis-server", "--requirepass", testRedisPassword),
	)
	if err != nil {
		t.Fatalf("fail to up container: %v", err)
	}

	defer func() {
		if err := redisContainer.Terminate(context.Background()); err != nil {
			fmt.Println("fail close container", err)
		}
	}()

	port, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("fail to get port: %v", err)
	}

	os.Setenv("REDIS_HOST", localHost)
	os.Setenv("REDIS_PORT", port.Port())
	os.Setenv("REDIS_PASSWORD", testRedisPassword)
	os.Setenv("REDIS_TIMEOUT", "30s")

	defer func() {
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("REDIS_PORT")
		os.Unsetenv("REDIS_PASSWORD")
		os.Unsetenv("REDIS_TIMEOUT")
	}()

	client, err := pkg_repository_redis.NewRedisClient(pkg_repository_redis.NewRedisConfigMust())
	if err != nil {
		t.Fatalf("fail to get client: %v", err)
	}

	cache := NewCouriersCache(client)

	t.Run("test redis", func(t *testing.T) {
		order := courier_domains.NewDispetchedOrder(1, 1, 1)

		if err := cache.SaveToCache(ctx, order); err != nil {
			t.Errorf("fail SaveToCache: %v", err)
		}

		gotOrder, err := cache.GetOrder(ctx, 1)
		if err != nil {
			t.Errorf("fail GetOrder: %v", err)
		}

		if gotOrder.OrderID != order.OrderID || gotOrder.Version != order.Version || gotOrder.CourierID != order.CourierID {
			t.Errorf("gotOrder is not order!")
		}

		if err := cache.DeleteOrder(ctx, 1); err != nil {
			t.Errorf("fail DeleteOrder: %v", err)
		}
	})
}
