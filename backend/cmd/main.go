// @title Simple Swagger Example
// @version 1.0
// @description Простейший пример API с Swagger
// @host localhost:8080
// @BasePath /

package main

import (
	"fmt"
	"log"
    "net/http"

    "github.com/gin-gonic/gin"
	"my_app_backend/internal/db"
    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    _ "my_app_backend/docs"
)

// Error handler func
func handleError(c *gin.Context, statusCode int, message string) {
    // Отправляем JSON с ошибкой
    c.JSON(statusCode, gin.H{
        "error": message,
    })
}

// @Summary Проверка сервера
// @ID ping
// @Description Тестовый эндпоинт для проверки работы сервера
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /ping [get]
func pingHandler(c *gin.Context) {
    simulatedError := false
    if simulatedError {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func main() {
	fmt.Println("🚀 Запуск приложения...")

	// Подключаемся к базе данных
	db.ConnectDB()

	// Создаём новый экземпляр Gin
	router := gin.Default()

	// Подключаем обработчик пинга
	router.GET("/ping", pingHandler)

	// Swagger UI доступен по адресу /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запускаем сервер
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
