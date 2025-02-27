.PHONY: up down logs migrate run build tidy

# Путь к docker-compose (запускаем из корня проекта)
DOCKER_COMPOSE = docker-compose -f ../docker/docker-compose.yml

# Запускаем контейнеры
up:
	$(DOCKER_COMPOSE) up --build -d

# Останавливаем контейнеры
down:
	$(DOCKER_COMPOSE) down

# Показываем логи бэкенда
logs:
	$(DOCKER_COMPOSE) logs -f backend

# Применяем миграции
migrate:
	docker exec backend_app migrate -path /app/migrations -database "postgres://postgres:password@postgres:5432/spaced_repetition_db?sslmode=disable" up

# Запускаем бэкенд без Docker
run:
	go run cmd/main.go

# Собираем Go-приложение
build:
	go build -o bin/app cmd/main.go

# Обновляем зависимости Go
tidy:
	go mod tidy
