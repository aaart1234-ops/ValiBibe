package controller

import (
    "net/http"

	"github.com/gin-gonic/gin"
	"my_app_backend/internal/service"
)

// AuthController отвечает за обработку запросов, связанных с авторизацией
type AuthController struct {
    authService service.AuthService
}

// NewAuthController создает новый экземпляр контроллера
func NewAuthController(authService service.AuthService) *AuthController {
    return &AuthController{authService: authService}
}

// RegisterUserHandler обрабатывает регистрацию пользователя
func (c *AuthController) RegisterUserHandler(ctx *gin.Context) {
    var request struct {
        Email string `json:"email"`
        Password string `json:"password"`
        Nickname string `json:"nickname"`
    }

    // Парсим JSON-запрос
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
        return
    }

    // Вызываем сервис для регистрации
    user, err := c.authService.RegisterUser(request.Email, request.Password, request.Nickname)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// Успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Пользователь зарегистрирован",
		"user":    user,
	})
}

// LoginUserHandler обрабатывает вход пользователя
func (c *AuthController) LoginUserHandler(ctx *gin.Context) {
    var request struct {
        Email string `json:"email"`
        Password string `json:"password"`
    }

    // Парсим JSON-запрос
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
        return
    }

    // Вызываем сервис для аутентификации
    token, err := c.authService.LoginUser(request.Email, request.Password)
    if err != nil {
    	ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    	return
    }

    // Успешный ответ с токеном
    ctx.JSON(http.StatusOK, gin.H{
        "message": "Вход выполнен успешно",
        "token":   token,
    })
}

















