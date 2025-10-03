package integration

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"valibibe/internal/models"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Note{}, &models.Folder{}, &models.Tag{}, &models.NoteTag{})
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}
