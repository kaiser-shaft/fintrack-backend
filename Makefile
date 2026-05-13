MAKEFLAGS += --no-print-directory
ENV_FILE ?= .env

-include $(ENV_FILE)
export

DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= password
DB_NAME ?= fintrack
DB_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: run
run:
	@CONFIG_PATH=$(ENV_FILE) go run ./cmd/api/main.go

.PHONY: up
env-up:
	@docker compose up -d db

.PHONY: env-down
env-down:
	@docker compose down --remove-orphans

.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then\
		echo "Ошибка: нужно указать имя миграции. Пример: make migrate-create name=init"; \
		exit 1; \
	fi; \
	migrate create -ext sql -dir ./migrations -seq $(name);

.PHONY: migrate
migrate:
	@if [ -z "$(action)" ]; then\
		echo "Ошибка: нужно указать action. Пример: make migrate action=\"up 2\""; \
		exit 1; \
	fi; \
	migrate -database "$(DB_URL)" -path=./migrations $(action)

.PHONY: migrate-up
migrate-up:
	@make migrate action="up"

.PHONY: migrate-down
migrate-down:
	@make migrate action="down"
