package integration

import (
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    "my_app_backend/internal/models"
    "my_app_backend/internal/repository"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to connect to test database: %v", err)
    }

    // Прогоняем миграции
    err = db.AutoMigrate(&models.User{})
    if err != nil {
        t.Fatalf("failed to migrate test database: %v", err)
    }

    return db
}

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