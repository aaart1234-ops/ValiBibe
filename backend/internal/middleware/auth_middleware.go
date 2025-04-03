package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"my_app_backend/internal/config"
)

// AuthMiddleware проверяет JWT-токен
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Извлекаем заголовок Authorization
        authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

        // Проверяем, что заголовок имеет формат "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
        }

        // Разбираем токен
        tokenString := parts[1]
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        			// Проверяем метод подписи
        			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        				return nil, fmt.Errorf("unexpected signing method")
        			}
        			return []byte(config.GetJWTSecret()), nil
        })

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

        // Извлекаем user_id из payload токена
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || claims["user_id"] == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
        }

        // Добавляем user_id в контекст запроса
        c.Set("user_id", claims["user_id"])
        c.Next()
    }
}