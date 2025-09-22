package bootstrap

import (
	"github.com/gin-gonic/gin"
	"valibibe/internal/db"
	"valibibe/internal/repository"
	"valibibe/internal/service"
	"valibibe/internal/controller"
	"valibibe/internal/router"
	"valibibe/internal/middleware"
	_ "valibibe/docs"
)

func InitializeApp() (*gin.Engine, error) {
	// Подключение к БД
	db.ConnectDB()
	database := db.GetDB()

	// Репозитории
	userRepo := repository.NewUserRepository(database)
	noteRepo := repository.NewNoteRepository(database)
	folderRepo := repository.NewFolderRepo(database)
	tagRepo := repository.NewTagRepository(database)

	// Сервисы
	tokenService := service.NewTokenService()
	authService := service.NewAuthService(userRepo, tokenService)
	noteService := service.NewNoteService(noteRepo)
	folderService := service.NewFolderService(folderRepo)
	tagService := service.NewTagService(tagRepo)

	// Контроллеры
	authController := controller.NewAuthController(authService)
	noteController := controller.NewNoteController(noteService)
	folderController := controller.NewFolderController(folderService)
	tagController := controller.NewTagController(tagService)

	// Инициализация Gin
	engine := gin.Default()
    engine.Use(middleware.CORSMiddleware())
	// Роутинг
	router.SetupRoutes(engine, tokenService, authController, noteController, folderController, tagController)

	return engine, nil
}
