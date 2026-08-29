//go:build integration

package dispetch_confirm_repository_postgres

import (
	"context"
	dispetch_domains "couriers/internal/Dispetch/core/domains"
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
	log              = "database system is ready to accept connections"
	localHost        = "127.0.0.1"
)

func TestDatabase_integration(t *testing.T) {
	ctx := t.Context()

	pgContainer, err := postgres.Run(
		ctx,
		postgresImage,
		postgres.WithDatabase(testDatabaseName),
		postgres.WithPassword(testPassword),
		postgres.WithUsername(testUser),
		testcontainers.WithWaitStrategy(wait.ForLog(log).WithOccurrence(2)),
	)
	if err != nil {
		t.Fatalf("fail to create container: %v", err)
	}

	defer func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			fmt.Println("fail to close container", err)
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

	CREATE TABLE app.orders (
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    name VARCHAR(20) NOT NULL CHECK(char_length(name) BETWEEN 3 AND 20),
    price BIGINT NOT NULL,
    is_complete BOOLEAN NOT NULL,
    city VARCHAR(100) NOT NULL CHECK(
        char_length(city) BETWEEN 1 AND 100
        AND
        city ~ '^[A-Z]'
    ),
    user_id INTEGER NOT NULL,
    courier_id INTEGER,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES app.users(id),
    CONSTRAINT fk_courier_id FOREIGN KEY (courier_id) REFERENCES app.couriers(id)
	);
	`

	contextWithTime, cancel := context.WithTimeout(ctx, pool.GetTimeot())
	defer cancel()

	if _, err := pool.Exec(contextWithTime, sqlRequestInit); err != nil {
		t.Fatalf("fail to create tables: %v", err)
	}

	sqlRequestInsert := `
	INSERT INTO app.couriers (login, password, city, orders_complete, is_free)
	VALUES ($1, $2, $3, $4, $5);
	`

	contextWithTimeIns, cancel := context.WithTimeout(ctx, pool.GetTimeot())
	defer cancel()

	if _, err := pool.Exec(contextWithTimeIns, sqlRequestInsert, "123456789", "123456789", "Moscow", 0, true); err != nil {
		t.Fatalf("fail to insert user: %v", err)
	}

	sqlRequestInsertUser := `
	INSERT INTO app.users (login, password, city)
	VALUES ($1, $2, $3);
	`

	contextWithTimeInsUser, cancel := context.WithTimeout(ctx, pool.GetTimeot())
	defer cancel()

	if _, err := pool.Exec(contextWithTimeInsUser, sqlRequestInsertUser, "123456789", "123456789", "Moscow"); err != nil {
		t.Fatalf("fail to insert user: %v", err)
	}

	sqlRequestInsertOrder := `
	INSERT INTO app.orders (name, price, is_complete, city, user_id)
	VALUES ($1, $2, $3, $4, $5);
	`

	contextWithTimeInsOrder, cancel := context.WithTimeout(ctx, pool.GetTimeot())
	defer cancel()

	if _, err := pool.Exec(contextWithTimeInsOrder, sqlRequestInsertOrder, "123456789", 1000, false, "Moscow", 1); err != nil {
		t.Fatalf("fail to insert user: %v", err)
	}

	db := NewConfirmPostgresDatabase(pool)

	t.Run("confirm", func(t *testing.T) {
		confirm := dispetch_domains.NewConfirm(1, 1)
		if err := db.ConfirmOrder(ctx, confirm); err != nil {
			t.Errorf("fail to confirm: %v", err)
		}
	})
}
