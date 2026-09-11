include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d cohesive-postgres

env-down:
	@docker compose down cohesive-postgres

redis-up:
	@docker compose up -d cohesive-redis

redis-down:
	@docker compose down -d cohesive-redis

env-cleanup:
	@read -p "Отчистить все файлы окружения? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down cohesive-postgres port-forwarder && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Файлы окружения отчищены."; \
	else \
		echo "Отмена."; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm cohesive-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm cohesive-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@cohesive-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

cohesive-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/cohesive/main.go

logs-cleanup:
	@read -p "Отчистить все логи? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Логи отчищены."; \
	else \
		echo "Отмена."; \
	fi

cohesive-deploy:
	@docker compose up -d --build cohesive

cohesive-undeploy:
	@docker compose down cohesive

ps:
	@docker compose ps

refresh-swag:
	@swag init -g main.go -o ../../docs -d ./cmd/cohesive,./internal/features,./internal/core --parseDependency --parseInternal
	