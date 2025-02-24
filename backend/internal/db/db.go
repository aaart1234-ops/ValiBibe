package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool" // Импорт библиотеки для работы с пулом соединений PostgreSQL
)

// DB — глобальная переменная для хранения пула соединений с БД.
var DB *pgxpool.Pool

// ConnectDB устанавливает соединение с базой данных.
func ConnectDB() {
	// Строка подключения (DSN — Data Source Name), содержит:
	// - логин (postgres)
	// - пароль (password)
	// - адрес хоста (localhost)
	// - порт (5432)
	// - имя базы данных (spaced_repetition_db)
	// - sslmode=disable (отключает SSL, используется для локальной разработки)
	dsn := "postgres://postgres:password@localhost:5432/spaced_repetition_db?sslmode=disable"

	// Парсим строку подключения в структуру конфигурации pgxpool
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Ошибка парсинга DSN: %v", err) // Завершаем выполнение программы, если произошла ошибка
	}

	// Создаём пул соединений с БД на основе конфигурации
	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err) // Если не удалось подключиться — завершаем программу
	}

	// Проверяем соединение с БД с помощью Ping
	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Ошибка проверки соединения: %v", err)
	}

	// Если все прошло успешно, выводим сообщение и сохраняем пул соединений в глобальную переменную DB
	fmt.Println("✅ Подключение к БД успешно!")
	DB = dbpool
}
