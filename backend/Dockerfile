# Используем базовый образ Go
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum из backend/
COPY backend/go.mod backend/go.sum ./

# Загружаем зависимости
RUN go mod tidy

# Копируем весь код проекта из backend/
COPY backend/ ./

# Компилируем приложение
RUN go build -o main cmd/main.go

# Новый минимальный образ для запуска
FROM ubuntu:latest
WORKDIR /app
COPY --from=builder /app/main .

# Запускаем приложение
CMD ["./main"]
