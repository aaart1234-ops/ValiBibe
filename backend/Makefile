.PHONY: up down logs migrate run build tidy

# Путь к docker-compose (запускаем из корня проекта)
DOCKER_COMPOSE = docker-compose -f ../docker/docker-compose.yml

# Запускаем контейнеры
up:
	$(DOCKER_COMPOSE) --env-file ../.env up --build -d

# Останавливаем контейнеры
down:
	$(DOCKER_COMPOSE) --env-file ../.env down

# Показываем логи бэкенда
logs:
	$(DOCKER_COMPOSE) logs -f backend

# Применяем миграции
migrate:
	docker run --rm --network=host \
		-v $(shell pwd)/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "postgres://postgres:password@localhost:5432/spaced_repetition_db?sslmode=disable" up

# Запускаем бэкенд без Docker
run:
	go run cmd/main.go

# Собираем Go-приложение
build:
	go build -o bin/app cmd/main.go

# Обновляем зависимости Go
tidy:
	go mod tidy

# Создать новую миграцию
create-migration:
	ifndef NAME
		$(error "Укажи имя миграции: make create-migration NAME=create_notes_table")

	endif
		mkdir -p migrations
		docker run --rm -v $(shell pwd)/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq $(NAME)

# Откатить одну миграцию вниз
migrate-down:
	docker run --rm --network=host \
		-v $(shell pwd)/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "postgres://postgres:password@localhost:5432/spaced_repetition_db?sslmode=disable" \
		down 1

# Проверить текущую версию миграции
migrate-status:
	docker run --rm --network=host \
		-v $(shell pwd)/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "postgres://postgres:password@localhost:5432/spaced_repetition_db?sslmode=disable" \
		version

# Принудительно установить версию миграции (если застряли в dirty-состоянии)
# Пример: make migrate-force VERSION=1
migrate-force:
	docker run --rm --network=host \
		-v $(shell pwd)/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "postgres://postgres:password@localhost:5432/spaced_repetition_db?sslmode=disable" \
		force $(VERSION)

# Полный откат всех миграций
migrate-reset:
	docker run --rm --network=host \
		-v $(shell pwd)/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "postgres://postgres:password@localhost:5432/spaced_repetition_db?sslmode=disable" \
		down
