package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"my_app_backend/internal/models"
)

// DB — глобальная переменная для хранения подключения GORM
var DB *gorm.DB

// ConnectDB устанавливает соединение с базой данных
func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("DSN:", dsn)
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL не задан!")
	}

	// Открываем соединение с базой через GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Включаем логирование SQL-запросов
	})
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}

	// Проверяем соединение
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Ошибка получения SQL соединения: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Ошибка проверки соединения: %v", err)
	}

	fmt.Println("✅ Подключение к БД успешно!")

	// Автоматически создаем таблицы
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("❌ Ошибка миграции: %v", err)
	}
	fmt.Println("✅ Миграции выполнены!")

	// Сохраняем подключение в глобальную переменную
	DB = db
}

// GetDB возвращает текущее подключение к БД
func GetDB() *gorm.DB {
    return DB
}