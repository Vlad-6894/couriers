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
	@docker compose run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@test-task-postgres:5432/${POSTGRES_DB}?sslmode=disable \
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