package utils

import (
    "github.com/google/uuid"
	"gorm.io/gorm"
)

// SetUUIDIfNil заполняет UUID, если он пустой
func SetUUIDIfNil(id *uuid.UUID) func(*gorm.DB) error {
    return func(tx *gorm.DB) error {
        if *id == uuid.Nil {
            *id = uuid.New()
        }
        return nil
    }
}

