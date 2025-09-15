package unit

import (
    "errors"
    "testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"

	"valibibe/internal/models"
	"valibibe/internal/service"
)

// Мок UserRepository
type MockUserRepo struct {
    mock.Mock
}

func (m *MockUserRepo) GetUserByEmail(email string) (*models.User, error) {
    args := m.Called(email)
    if user := args.Get(0); user != nil {
        return user.(*models.User), args.Error(1)
    }
    return nil, args.Error(1)
}

func (m *MockUserRepo) GetUserByID(id string) (*models.User, error) {
	args := m.Called(id)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}

func (m *MockUserRepo) CreateUser(user *models.User) error {
    args := m.Called(user)
    return args.Error(0)
}

// Мок TokenService
type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateToken(userID uuid.UUID) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) ValidateToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	if token := args.Get(0); token != nil {
		return token.(*jwt.Token), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockToken := new(MockTokenService)

	authService := service.NewAuthService(mockRepo, mockToken)

	email := "test@example.com"
	nickname := "TestUser"
	password := "securepassword"

	mockRepo.On("GetUserByEmail", email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("CreateUser", mock.Anything).Return(nil)

	user, err := authService.RegisterUser(email, password, nickname)

	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, nickname, user.Nickname)
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockToken := new(MockTokenService)
	authService := service.NewAuthService(mockRepo, mockToken)

	email := "test@example.com"
	password := "securepassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	userID := uuid.New()

	mockUser := &models.User{
		ID:           userID,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	mockRepo.On("GetUserByEmail", email).Return(mockUser, nil)
	mockToken.On("GenerateToken", userID).Return("mocked.jwt.token", nil)

	token, err := authService.LoginUser(email, password)

	assert.NoError(t, err)
	assert.Equal(t, "mocked.jwt.token", token)
	mockRepo.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}

func TestLoginUser_WrongEmail(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockToken := new(MockTokenService)
	authService := service.NewAuthService(mockRepo, mockToken)

	email := "wrong@example.com"
	password := "password"

	mockRepo.On("GetUserByEmail", email).Return(nil, errors.New("not found"))

	token, err := authService.LoginUser(email, password)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.EqualError(t, err, "пользователь не найден")
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockToken := new(MockTokenService)
	authService := service.NewAuthService(mockRepo, mockToken)

	email := "test@example.com"
	password := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	mockUser := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	mockRepo.On("GetUserByEmail", email).Return(mockUser, nil)

	token, err := authService.LoginUser(email, password)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.EqualError(t, err, "неверный пароль")
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
    mockRepo := new(MockUserRepo)
    mockToken := new(MockTokenService)
    authService := service.NewAuthService(mockRepo, mockToken)

    userID := uuid.New().String()
    expectedUser := &models.User{
        ID:       uuid.MustParse(userID),
        Email:    "test@example.com",
        Nickname: "Tester",
    }

    mockRepo.On("GetUserByID", userID).Return(expectedUser, nil)

    user, err := authService.GetUserByID(userID)

    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, expectedUser.Email, user.Email)
    assert.Equal(t, expectedUser.Nickname, user.Nickname)

    mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockToken := new(MockTokenService)
	authService := service.NewAuthService(mockRepo, mockToken)

	userID := uuid.New().String()

	// Предположим, репозиторий вернет gorm.ErrRecordNotFound
	mockRepo.On("GetUserByID", userID).Return(nil, gorm.ErrRecordNotFound)

	user, err := authService.GetUserByID(userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockRepo.AssertExpectations(t)
}

