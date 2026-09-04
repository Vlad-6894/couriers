include .env
export

export PROJECT_ROOT=$(shell pwd)

start-postgres:
	@docker compose up -d couriers-postgres-db

finish-postgres:
	@docker compose down couriers-postgres-db

cleanup-postgres:
	@read -p "Очистить pg_data? ВНИМАНИЕ! Опасность утери данных! [y/N]: " choice; \
	if [ "$$choice" = "y" ] || [ "$$choice" = "Y" ]; then \
		docker compose down couriers-postgres-db couriers-postgres-port-forwarder && \
		sudo rm -rf ${PROJECT_ROOT}/out/pg_data && \
		echo "Очищено"; \
	else \
		echo "Операция отменена"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Вы не передали название миграции!" \
		exit 1; \
	fi;
	@docker compose run --rm couriers-postgres-db-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Нет параметра action"; \
		exit 1; \
	fi;
	@docker compose run --rm couriers-postgres-db-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@couriers-postgres-db:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

start-port-forwarder:
	@docker compose up -d couriers-postgres-port-forwarder

finish-port-forwarder:
	@docker compose down couriers-postgres-port-forwarder

start-redis:
	@docker compose up -d couriers-redis

finish-redis:
	@docker compose down couriers-redis

start-kafka:
	@docker compose up -d couriers-kafka

finish-kafka:
	@docker compose down couriers-kafka

print-jwt:
	@openssl rand -hex 32

start-auth-service:
	@docker compose up -d auth

finish-auth-service:
	@docker compose down auth

start-orders-service:
	@docker compose up -d orders

finish-orders-service:
	@docker compose down orders

start-dispetch-service:
	@docker compose up -d dispetch

finish-dispetch-service:
	@docker compose down dispetch

start-couriers-service:
	@docker compose up -d couriers

finish-couriers-service:
	@docker compose down couriers

generate-mocks:
	@go generate ./...

start-unit-tests-auth-service:
	@go test -v ./internal/Auth/features/service

start-integration-tests-auth-service:
	@go test -v -tags=integration ./internal/Auth/features/repository/postgres

start-unit-tests-orders-service:
	@go test -v ./internal/Orders/features/service

start-integration-tests-orders-service:
	@go test -v -tags=integration ./internal/Orders/features/repository/...

start-unit-tests-dispetch-service-orders:
	@go test -v ./internal/Dispetch/features/order/service

start-unit-tests-dispetch-service-confirm:
	@go test -v ./internal/Dispetch/features/confirm/service

start-integration-tests-dispetch-service-orders-broker:
	@go test -v -tags=integration ./internal/Dispetch/features/order/repository/kafka

start-integration-tests-dispetch-service-orders-database:
	@go test -v -tags=integration ./internal/Dispetch/features/order/repository/postgres

start-integration-tests-dispetch-service-orders-cache:
	@go test -v -tags=integration ./internal/Dispetch/features/order/repository/redis

start-integration-tests-dispetch-service-confirm-database:
	@go test -v -tags=integration ./internal/Dispetch/features/confirm/repository/postgres

start-integration-tests-dispetch-service-confirm:
	@go test -v -tags=integration ./internal/Dispetch/features/confirm/repository/...

start-unit-tests-couriers-service:
	@go test -v ./internal/Couriers/features/service

start-integration-tests-couriers-service-cache:
	@go test -v -tags=integration ./internal/Couriers/features/repository/redis

start-integration-tests-couriers-service-broker:
	@go test -v -tags=integration ./internal/Couriers/features/repository/kafka