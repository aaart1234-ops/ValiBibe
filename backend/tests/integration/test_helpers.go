package integration

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"my_app_backend/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Note{})
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}
