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

func main() {
	fmt.Println("🚀 Запуск приложения...")

	// Подключаемся к базе данных
	db.ConnectDB()

	// Create new instance Gin
	router := gin.Default()

	// Add test endpoint
	router.GET("/ping", func(c *gin.Context) {
	    // Симулируем ошибку (например, если что-то пошло не так)
	        simulatedError := false
	        if simulatedError {
	            handleError(c, http.StatusInternalServerError, "Ошибка сервера")
	            return // Выходим из обработчика, чтобы не продолжать выполнение кода
	        }

            // Если ошибки нет, отправляем стандартный ответ
    		c.JSON(http.StatusOK, gin.H{"message": "pong"})
    })

    // Swagger UI доступен по адресу /swagger/index.html
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Launch server on 8080 port
    err := router.Run(":8080");
    if err != nil {
        log.Fatal("Ошибка запуска сервера: ", err);
    }
}
