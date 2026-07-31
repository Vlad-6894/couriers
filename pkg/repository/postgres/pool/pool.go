package pkg_postgres_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	GetTimeot() time.Duration
}

type PostgresConnectionPool struct {
	*pgxpool.Pool
	timeout time.Duration
}

func NewPostgresConnectionPool(ctx context.Context, config PostgresConnectionPoolConfig) (*PostgresConnectionPool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxconfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgxconfig error: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool error: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping pgxpool error: %w", err)
	}

	return &PostgresConnectionPool{
		Pool:    pool,
		timeout: config.Timeout,
	}, nil
}

func (p *PostgresConnectionPool) GetTimeot() time.Duration {
	return p.timeout
}
