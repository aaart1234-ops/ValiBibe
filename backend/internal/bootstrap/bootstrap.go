package bootstrap

import (
	_ "valibibe/docs"
	"valibibe/internal/controller"
	"valibibe/internal/db"
	"valibibe/internal/middleware"
	"valibibe/internal/repository"
	"valibibe/internal/router"
	"valibibe/internal/service"

	"github.com/gin-gonic/gin"
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
	assignFolderService := service.NewAssignFolderService(noteRepo, folderRepo)
	folderService := service.NewFolderService(folderRepo)
	tagService := service.NewTagService(tagRepo)
	noteTagService := service.NewNoteTagService(noteRepo, tagRepo)
	reviewSessionService := service.NewReviewSessionService(noteRepo)

	// Контроллеры
	authController := controller.NewAuthController(authService)
	noteController := controller.NewNoteController(noteService, assignFolderService)
	folderController := controller.NewFolderController(folderService)
	tagController := controller.NewTagController(tagService)
	noteTagController := controller.NewNoteTagController(noteTagService)
	reviewSessionController := controller.NewReviewSessionController(reviewSessionService)

	// Инициализация Gin
	engine := gin.Default()
	engine.Use(middleware.CORSMiddleware())
	// Роутинг
	router.SetupRoutes(engine, tokenService, authController, noteController, folderController, tagController, noteTagController, reviewSessionController)

	return engine, nil
}
