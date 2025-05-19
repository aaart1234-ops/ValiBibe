// @title Simple Swagger Example
// @version 1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description backend API для проекта
// @host localhost:8080
// @BasePath /

package main

import (
	"fmt"
	"log"
	"my_app_backend/internal/bootstrap"
)

func main() {
	fmt.Println("🚀 Запуск приложения...")

	app, err := bootstrap.InitializeApp()
	if err != nil {
		log.Fatal("❌ Ошибка инициализации: ", err)
	}

	if err := app.Run(":8181"); err != nil {
		log.Fatal("❌ Ошибка запуска сервера: ", err)
	}
}
