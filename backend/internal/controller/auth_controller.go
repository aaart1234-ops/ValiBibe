package controller

import (
    "fmt"
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
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя по email, паролю и нику
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.RegisterRequest true "Данные для регистрации"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (c *AuthController) RegisterUserHandler(ctx *gin.Context) {
    var request struct {
        Email string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
        Nickname string `json:"nickname" binding:"required"`
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
// @Summary Вход пользователя
// @Description Аутентификация пользователя по email и паролю, возвращает JWT-токен
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.LoginRequest true "Данные для входа"
// @Success 200 {object} map[string]interface{} "Успешный вход и токен"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 401 {object} map[string]string "Неверный email или пароль"
// @Router /auth/login [post]
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

// MeHandler возвращает информацию о текущем пользователе
// @Summary Получить текущего пользователя
// @Description Возвращает данные пользователя на основе JWT-токена
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/me [get]
func (c *AuthController) MeHandler(ctx *gin.Context) {
    fmt.Println("🔍 MeHandler вызывается") // 👈 сюда

    userID, exists := ctx.Get("user_id")
    if !exists {
        fmt.Println("⛔️ user_id не найден в context")
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
        return
    }

    fmt.Println("✅ user_id:", userID)

    user, err := c.authService.GetUserByID(userID.(string))
    if err != nil {
        fmt.Println("⛔️ Ошибка при получении пользователя:", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить пользователя"})
        return
    }

    fmt.Println("📤 Отправляем данные пользователя:", user)
    ctx.JSON(http.StatusOK, user)
}

// LogoutHandler обрабатывает выход пользователя
// @Summary Выход из системы
// @Description Инвалидирует токен на клиенте (сервер токены не хранит)
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
// LogoutHandler обрабатывает выход пользователя
func (ac *AuthController) LogoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "✅ Вы вышли из системы"})
}




















