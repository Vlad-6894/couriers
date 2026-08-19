package main

import (
	"context"
	pkg_logger "couriers/pkg/logger"
	pkg_postgres_pool "couriers/pkg/repository/postgres/pool"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := pkg_logger.NewLogger(pkg_logger.NewLoggerConfigMust())
	if err != nil {
		fmt.Println("init logger error!")
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("init connection pool!")
	pool, err := pkg_postgres_pool.NewPostgresConnectionPool(ctx, pkg_postgres_pool.NewPostgresConnectionConfigMust())
	if err != nil {
		logger.Fatal("connection pool init error!")
	}
	defer pool.Close()
}
