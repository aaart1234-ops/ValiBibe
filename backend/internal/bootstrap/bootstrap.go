package bootstrap

import (
	"github.com/gin-gonic/gin"
	"my_app_backend/internal/db"
	"my_app_backend/internal/repository"
	"my_app_backend/internal/service"
	"my_app_backend/internal/controller"
	"my_app_backend/internal/router"
	_ "my_app_backend/docs"
)

func InitializeApp() (*gin.Engine, error) {
	// Подключение к БД
	db.ConnectDB()
	database := db.GetDB()

	// Репозитории
	userRepo := repository.NewUserRepository(database)

	// Сервисы
	tokenService := service.NewTokenService()
	authService := service.NewAuthService(userRepo, tokenService, database)

	// Контроллеры
	authController := controller.NewAuthController(authService)

	// Инициализация Gin
	engine := gin.Default()

	// Роутинг
	router.SetupRoutes(engine, tokenService, authController)

	return engine, nil
}
