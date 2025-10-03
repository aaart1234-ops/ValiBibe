package integration

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"valibibe/internal/repository"
	"valibibe/internal/service"
)

func TestAuthService_FullFlow(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Ошибка при загрузке .env файла")
	}

	db := SetupTestDB(t)

	userRepo := repository.NewUserRepository(db)
	tokenService := service.NewTokenService()
	authService := service.NewAuthService(userRepo, tokenService)

	email := "testuser@example.com"
	password := "securePassword"
	nickname := "TestUser"

	// 1. Register
	createdUser, err := authService.RegisterUser(email, password, nickname)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, email, createdUser.Email)
	assert.Equal(t, nickname, createdUser.Nickname)
	assert.NotEmpty(t, createdUser.ID)

	// 2. Login
	token, err := authService.LoginUser(email, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// 3. GetUserByID
	fetchedUser, err := authService.GetUserByID(createdUser.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, fetchedUser.ID)
	assert.Equal(t, createdUser.Email, fetchedUser.Email)
}
