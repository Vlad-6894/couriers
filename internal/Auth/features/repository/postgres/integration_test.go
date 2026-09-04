//go:build integration

package auth_repository_postgres

import (
	"context"
	auth_domains "couriers/internal/Auth/core/domains"
	pkg_postgres_pool "couriers/pkg/repository/postgres/pool"
	"fmt"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	postgresImage    = "postgres:18.1-bookworm"
	testDatabaseName = "auth_integration_db"
	testUser         = "test_user"
	testPassword     = "test123"
	localHost        = "127.0.0.1"
)

func TestAuthRepository_Integration(t *testing.T) {
	ctx := t.Context()

	pgContainer, err := postgres.Run(
		ctx,
		postgresImage,
		postgres.WithDatabase(testDatabaseName),
		postgres.WithUsername(testUser),
		postgres.WithPassword(testPassword),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
	)
	if err != nil {
		t.Fatalf("fail to start test container: %v", err)
	}

	defer func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			fmt.Println("fail close container: ", err)
		}
	}()

	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get port: %v", err)
	}

	os.Setenv("POSTGRES_HOST", localHost)
	os.Setenv("POSTGRES_PORT", mappedPort.Port())
	os.Setenv("POSTGRES_USER", testUser)
	os.Setenv("POSTGRES_PASSWORD", testPassword)
	os.Setenv("POSTGRES_DB", testDatabaseName)
	os.Setenv("POSTGRES_TIMEOUT", "30s")

	defer func() {
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("POSTGRES_DB")
		os.Unsetenv("POSTGRES_TIMEOUT")
	}()

	pool, err := pkg_postgres_pool.NewPostgresConnectionPool(
		context.Background(),
		pkg_postgres_pool.NewPostgresConnectionConfigMust(),
	)
	if err != nil {
		t.Fatalf("fail to get pool: %v, port: %v", err, mappedPort.Port())
	}
	defer pool.Close()

	sqlRequestInit := `
	CREATE SCHEMA IF NOT EXISTS app;

	CREATE TABLE IF NOT EXISTS app.users (
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    login VARCHAR(100) NOT NULL CHECK(char_length(login) BETWEEN 8 AND 100),
    password VARCHAR(100) NOT NULL CHECK(char_length(password) BETWEEN 8 AND 100),
    city VARCHAR(100) NOT NULL CHECK(
        char_length(city) BETWEEN 1 AND 100
        AND
        city ~ '^[A-Z]'
    ),

    UNIQUE(login)
	);

	CREATE TABLE IF NOT EXISTS app.couriers (
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    login VARCHAR(100) NOT NULL CHECK(char_length(login) BETWEEN 8 AND 100),
    password VARCHAR(100) NOT NULL CHECK(char_length(password) BETWEEN 8 AND 100),
    city VARCHAR(100) NOT NULL CHECK(
        char_length(city) BETWEEN 1 AND 100
        AND
        city ~ '^[A-Z]'
    ),
    orders_complete BIGINT NOT NULL,
    is_free BOOLEAN NOT NULL,

    UNIQUE(login)
	);
	`

	contextWithTime, cancel := context.WithTimeout(ctx, pool.GetTimeot())
	defer cancel()

	if _, err := pool.Exec(contextWithTime, sqlRequestInit); err != nil {
		t.Fatalf("fail to create tables: %v", err)
	}

	authRepository := NewAuthRepository(pool)

	t.Run("Register and get user", func(t *testing.T) {
		user := auth_domains.NewRegUser("123456789", "123456789", "Moscow")

		user, err := authRepository.RegisterUser(ctx, user)
		if err != nil {
			t.Errorf("fail to register user: %v", err)
		}

		if _, err := authRepository.GetUser(ctx, user.Login); err != nil {
			t.Errorf("fail to get user: %v", err)
		}

		if _, err := authRepository.RegisterUser(ctx, user); err == nil {
			t.Errorf("fail to unique in databse")
		}
	})

	t.Run("Register and get courier", func(t *testing.T) {
		courier := auth_domains.NewRegCourier("123456789", "123456789", "Moscow")

		courier, err := authRepository.RegisterCourier(ctx, courier)
		if err != nil {
			t.Errorf("fail to register courier: %v", err)
		}

		if _, err := authRepository.GetCourier(ctx, courier.Login); err != nil {
			t.Errorf("fail to get courier: %v", err)
		}

		if _, err := authRepository.RegisterCourier(ctx, courier); err == nil {
			t.Errorf("fail to unique in databse")
		}
	})
}
