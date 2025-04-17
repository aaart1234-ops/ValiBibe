package service_tests

import (
    "os"
    "testing"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"my_app_backend/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestTokenService_GenerateAndValidateToken(t *testing.T) {
    // Устанавливаем переменную окружения
    os.Setenv("JWT_SECRET", "test_secret")

    tokenSvc := service.NewTokenService()

    // Создаем UUID пользователя
    userID := uuid.New()

    // Генерируем токен
    tokenString, err := tokenSvc.GenerateToken(userID)
    assert.NoError(t, err)
    assert.NotEmpty(t, tokenString)

    // Валидируем токен
    token, err := tokenSvc.ValidateToken(tokenString)
    assert.NoError(t, err)
    assert.NotNil(t, token)

    // Извлекаем claims
    claims, ok := token.Claims.(jwt.MapClaims)
    assert.True(t, ok)
    assert.Equal(t, userID.String(), claims["user_id"])
}