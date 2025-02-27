package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB — глобальная переменная для хранения пула соединений с БД.
var DB *pgxpool.Pool

// ConnectDB устанавливает соединение с базой данных.
func ConnectDB() {
	// Берем URL базы из переменной окружения
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL не задан!")
	}

	// Парсим строку подключения
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("❌ Ошибка парсинга DSN: %v", err)
	}

	// Создаём пул соединений с БД
	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}

	// Проверяем соединение
	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("❌ Ошибка проверки соединения: %v", err)
	}

	fmt.Println("✅ Подключение к БД успешно!")
	DB = dbpool
}
