package service

import (
    "errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"my_app_backend/internal/models"
	"my_app_backend/internal/repository"
)

// AuthService - интерфейс сервиса аутентификации
type AuthService interface {
    RegisterUser(email, password, nickname string) (*models.User, error)
    LoginUser(email, password string) (string, error)
    GetUserByID(userID string) (*models.User, error)
}

// authService - реализация AuthService
type authService struct {
    userRepo repository.UserRepository
    tokenService TokenService
}

// NewAuthService создает новый сервис аутентификации
func NewAuthService(userRepo repository.UserRepository, tokenService TokenService) AuthService {
    return &authService{userRepo: userRepo, tokenService: tokenService}
}

// RegisterUser регистрирует нового пользователя
func (s *authService) RegisterUser(email, password, nickname string) (*models.User, error) {
	// Проверяем, существует ли уже пользователь с таким email
	existingUser, err := s.userRepo.GetUserByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя
	user := &models.User{
		Email:        email,
		Nickname:     nickname,
		PasswordHash: string(hashedPassword),
	}

	// Сохраняем пользователя в базе
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// LoginUser выполняет аутентификацию пользователя и возвращает JWT-токен
func (s *authService) LoginUser(email, password string) (string, error) {
	// Ищем пользователя по email
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("пользователь не найден")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("неверный пароль")
	}

	// Генерируем JWT-токен
	token, err := s.tokenService.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Получает пользователя по ID
func (s *authService) GetUserByID(userID string) (*models.User, error) {
    return s.userRepo.GetUserByID(userID)
}