//go:build integration

package dispetch_redis_repository

import (
	"context"
	dispetch_domains "couriers/internal/Dispetch/core/domains"
	pkg_repository_redis "couriers/pkg/repository/redis"
	"fmt"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	tc_redis "github.com/testcontainers/testcontainers-go/modules/redis"
)

var (
	keyPrefix         = "couriers:pool"
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

	cache := NewRedisRepository(client)

	t.Run("test redis", func(t *testing.T) {
		firstDates := map[string]map[int]dispetch_domains.FreeCourierInfo{
			"Moscow": {
				1: dispetch_domains.NewFreeCourierInfo(1, 1),
				2: dispetch_domains.NewFreeCourierInfo(2, 1),
			},

			"London": {
				1: dispetch_domains.NewFreeCourierInfo(3, 1),
				2: dispetch_domains.NewFreeCourierInfo(4, 1),
			},
		}
		if err := cache.UpdateCache(ctx, firstDates); err != nil {
			t.Errorf("fail update cache: %v", err)
		}

		isUnique, err := cache.CheckUnique(ctx, 1)
		if err != nil || isUnique == false {
			t.Errorf("fail check unique: %v", err)
		}

		if _, _, err := cache.SearchCourier(ctx, "Moscow"); err != nil {
			t.Errorf("fail search courier: %v", err)
		}

		if err := cache.UpdateCache(ctx, firstDates); err != nil {
			t.Errorf("fail update cache: %v", err)
		}
	})
}
