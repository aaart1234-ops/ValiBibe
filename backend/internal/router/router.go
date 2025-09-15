package router

import (
    "time"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"valibibe/internal/controller"
	"valibibe/internal/middleware"
	"valibibe/internal/service"
)

func SetupRoutes(r *gin.Engine,
                    tokenService service.TokenService,
                    authController *controller.AuthController,
                    noteController *controller.NoteController,
                    folderController *controller.FolderController,
                ) {
	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Ping
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

    r.GET("/", func(c *gin.Context) {
    	c.Header("Content-Type", "text/html; charset=utf-8")
    	c.String(200, `
    		<!DOCTYPE html>
    		<html lang="ru">
    		<head>
    			<meta charset="UTF-8">
    			<title>My App Backend</title>
    			<style>
    				body {
    					font-family: Arial, sans-serif;
    					background-color: #f9f9f9;
    					padding: 40px;
    					color: #333;
    				}
    				h1 {
    					color: #007acc;
    				}
    				a {
    					color: #007acc;
    					text-decoration: none;
    				}
    			</style>
    		</head>
    		<body>
    			<h1>🚀 Приложение успешно запущено</h1>
    			<p><strong>Название:</strong> My App Backend</p>
    			<p><strong>Версия:</strong> 1.0</p>
    			<p><strong>Swagger UI:</strong> <a href="/swagger/index.html" target="_blank">Открыть документацию</a></p>
    			<p><strong>Время:</strong> ` + time.Now().Format("02 Jan 2006 15:04:05") + `</p>
    		</body>
    		</html>
    	`)
    })


	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.RegisterUserHandler)
		auth.POST("/login", authController.LoginUserHandler)
		auth.GET("/me", middleware.AuthMiddleware(tokenService), authController.MeHandler)
		auth.POST("/logout", middleware.AuthMiddleware(tokenService), authController.LogoutHandler)
	}

    // Notes
    notes := r.Group("/notes")
    notes.Use(middleware.AuthMiddleware(tokenService))
    {
        notes.POST("", noteController.CreateNote)
        notes.GET("", noteController.GetAllNotes)
        notes.GET("/:id", noteController.GetNoteByID)
        notes.PUT("/:id", noteController.UpdateNote)
        notes.DELETE("/:id", noteController.DeleteNote)
        notes.POST("/:id/archive", noteController.ArchiveNote)
        notes.POST("/:id/unarchive", noteController.UnArchiveNote)
        notes.POST("/:id/review", noteController.ReviewNoteHandler)
    }

    // Folders
    folders := r.Group("/folders")
    folders.Use(middleware.AuthMiddleware(tokenService))
    {
        folders.POST("", folderController.CreateFolder)
        folders.GET("/tree", folderController.GetFolderTree)
        folders.PUT("/:id", folderController.UpdateFolder)
        folders.DELETE("/:id", folderController.DeleteFolder)
    }
}
