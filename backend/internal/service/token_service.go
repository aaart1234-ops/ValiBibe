package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"errors"
	"os"
)

// TokenService определяет методы для работы с JWT
type TokenService interface {
	GenerateToken(userID uuid.UUID) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

// tokenService - реализация TokenService
type tokenService struct {
	secretKey string
}

// NewTokenService создает новый экземпляр сервиса токенов
func NewTokenService() TokenService {
	// Получаем секретный ключ из переменных окружения
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		panic("JWT_SECRET не задан в переменных окружения")
	}

	return &tokenService{secretKey: secretKey}
}

// GenerateToken создает JWT-токен для пользователя
func (s *tokenService) GenerateToken(userID uuid.UUID) (string, error) {
	// Создаем claims (полезные данные внутри токена)
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен действует 24 часа
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken проверяет и парсит токен
func (s *tokenService) ValidateToken(tokenString string) (*jwt.Token, error) {
	// Разбираем и валидируем токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что используется ожидаемый метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи токена")
		}
		return []byte(s.secretKey), nil
	})

	return token, err
}
