package middleware

import (
    "fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"valibibe/internal/service"
)

// AuthMiddleware принимает TokenService и проверяет JWT
func AuthMiddleware(tokenService service.TokenService) gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("🛡 Вошли в AuthMiddleware")

        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            fmt.Println("⛔️ Нет заголовка Authorization")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            fmt.Println("⛔️ Неправильный формат заголовка Authorization")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
            c.Abort()
            return
        }

        tokenString := parts[1]
        token, err := tokenService.ValidateToken(tokenString)
        if err != nil || !token.Valid {
            fmt.Println("⛔️ Невалидный токен:", err)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || claims["user_id"] == nil {
            fmt.Println("⛔️ Claims невалидные")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        userIDStr, ok := claims["user_id"].(string)
        if !ok {
            fmt.Println("⛔️ user_id не строка в токене")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in token"})
            c.Abort()
            return
        }

        fmt.Println("✅ Авторизация прошла, user_id:", userIDStr)
        c.Set("user_id", userIDStr)
        c.Next()
    }
}
