.PHONY: \
	default \
	up down logs \
	migrate migrate-down migrate-status migrate-force migrate-reset create-migration \
	run build tidy test \
	docs clean help

### -----------------------------------------
### Переменные
### -----------------------------------------

DB_URL = postgres://postgres:password@localhost:5432/spaced_repetition_db?sslmode=disable
MIGRATIONS_PATH = $(shell pwd)/migrations
DOCKER_COMPOSE = docker-compose -f ../docker/docker-compose.yml

### -----------------------------------------
### Команда по умолчанию
### -----------------------------------------

default: run  ## По умолчанию запускает backend локально

### -----------------------------------------
### Docker / dev-среда
### -----------------------------------------

up:            ## Запускает контейнеры
	$(DOCKER_COMPOSE) --env-file ../.env up --build -d

down:          ## Останавливает контейнеры
	$(DOCKER_COMPOSE) --env-file ../.env down

logs:          ## Показывает логи backend-контейнера
	$(DOCKER_COMPOSE) logs -f backend

### -----------------------------------------
### Миграции (через Docker образ migrate/migrate)
### -----------------------------------------

migrate:       ## Применяет все миграции
	docker run --rm --network=host \
		-v $(MIGRATIONS_PATH):/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$(DB_URL)" \
		up

create-migration: ## Создает новую миграцию. NAME=create_notes_table
ifndef NAME
	$(error Укажи имя миграции: make create-migration NAME=create_notes_table)
endif
	mkdir -p $(MIGRATIONS_PATH)
	docker run --rm \
		-v $(MIGRATIONS_PATH):/migrations \
		migrate/migrate \
		create -ext sql -dir /migrations -seq $(NAME)

migrate-down:  ## Откатывает одну миграцию
	docker run --rm --network=host \
		-v $(MIGRATIONS_PATH):/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$(DB_URL)" \
		down 1

migrate-status: ## Показывает текущую версию миграций
	docker run --rm --network=host \
		-v $(MIGRATIONS_PATH):/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$(DB_URL)" \
		version

migrate-force: ## Принудительно устанавливает версию миграции. VERSION=1
	docker run --rm --network=host \
		-v $(MIGRATIONS_PATH):/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$(DB_URL)" \
		force $(VERSION)

migrate-reset: ## Полный откат всех миграций
	docker run --rm --network=host \
		-v $(MIGRATIONS_PATH):/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$(DB_URL)" \
		down

### -----------------------------------------
### Go-команды (локальная разработка)
### -----------------------------------------

run:           ## Запускает backend без Docker
	go run cmd/main.go

build:         ## Собирает Go-бинарник
	go build -o bin/app cmd/main.go

tidy:          ## Обновляет зависимости Go
	go mod tidy

test:          ## Запускает все Go-тесты
	go test -v ./tests/...

### -----------------------------------------
### Документация
### -----------------------------------------

docs:          ## Генерация Swagger-доков
	@echo "Generating Swagger docs..."
	@swag init -g cmd/main.go -o docs

### -----------------------------------------
### Утилиты
### -----------------------------------------

clean:         ## Удаляет артефакты сборки
	rm -rf bin/*

help:          ## Показывает список команд и описания
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?##' Makefile | awk 'BEGIN {FS = ":.*?##"}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
