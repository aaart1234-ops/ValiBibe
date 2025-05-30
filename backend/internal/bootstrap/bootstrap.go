package bootstrap

import (
	"github.com/gin-gonic/gin"
	"my_app_backend/internal/db"
	"my_app_backend/internal/repository"
	"my_app_backend/internal/service"
	"my_app_backend/internal/controller"
	"my_app_backend/internal/router"
	"my_app_backend/internal/middleware"
	_ "my_app_backend/docs"
)

func InitializeApp() (*gin.Engine, error) {
	// Подключение к БД
	db.ConnectDB()
	database := db.GetDB()

	// Репозитории
	userRepo := repository.NewUserRepository(database)
	noteRepo := repository.NewNoteRepository(database)

	// Сервисы
	tokenService := service.NewTokenService()
	authService := service.NewAuthService(userRepo, tokenService)
	noteService := service.NewNoteService(noteRepo)

	// Контроллеры
	authController := controller.NewAuthController(authService)
	noteController := controller.NewNoteController(noteService)

	// Инициализация Gin
	engine := gin.Default()
    engine.Use(middleware.CORSMiddleware())
	// Роутинг
	router.SetupRoutes(engine, tokenService, authController, noteController)

	return engine, nil
}
