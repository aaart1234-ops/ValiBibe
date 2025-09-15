package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"valibibe/internal/middleware"
	"valibibe/internal/repository"
	"valibibe/internal/service"
)

func TestAuthMiddleware(t *testing.T) {
	db := setupTestDB(t)

	userRepo := repository.NewUserRepository(db)
	tokenService := service.NewTokenService()
	authService := service.NewAuthService(userRepo, tokenService)

	// создаем пользователя и токен
	user, _ := authService.RegisterUser("middleware@example.com", "password123", "MiddlewareUser")
	token, _ := authService.LoginUser("middleware@example.com", "password123")

	// создаем роутер с защищенным маршрутом
	router := gin.Default()
	router.GET("/protected", middleware.AuthMiddleware(tokenService), func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "access granted", "user_id": userID})
	})

	t.Run("Missing token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Invalid token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Valid token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var body map[string]interface{}
		_ = json.Unmarshal(resp.Body.Bytes(), &body)
		assert.Equal(t, "access granted", body["message"])
		assert.Equal(t, user.ID.String(), body["user_id"])
	})
}
