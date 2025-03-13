package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB — глобальная переменная для GORM
var DB *gorm.DB

// ConnectDB устанавливает соединение с базой данных через GORM
func ConnectDB() {
	// Берем URL базы из переменной окружения
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL не задан!")
	}

	// Открываем соединение с БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // Логируем только предупреждения и ошибки
	})
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}

	// Проверяем соединение
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Ошибка получения sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Ошибка проверки соединения: %v", err)
	}

	fmt.Println("✅ Подключение к БД через GORM успешно!")
	DB = db
}
