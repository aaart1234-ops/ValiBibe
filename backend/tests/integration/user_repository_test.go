package integration

import (
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"

    "valibibe/internal/models"
    "valibibe/internal/repository"
)

func TestUserRepository_CreateAndGetUserByID(t *testing.T) {
    db := setupTestDB(t)
    userRepo := repository.NewUserRepository(db)

    // Создаём нового пользователя
    user := &models.User{
        ID: uuid.New(),
        Email: "integration@example.com",
        Nickname: "IntegrationTest",
        PasswordHash: "hashedpassword123",
    }

    err := userRepo.CreateUser(user)
    assert.NoError(t, err)

    // Ищем по ID
    foundUser, err := userRepo.GetUserByID(user.ID.String())
    assert.NoError(t, err)
    assert.NotNil(t, foundUser)
    assert.Equal(t, user.Email, foundUser.Email)
    assert.Equal(t, user.Nickname, foundUser.Nickname)
}