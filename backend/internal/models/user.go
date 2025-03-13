package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// User представляет модель пользователя в базе данных.
type User struct {
    ID                uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    Nickname          string         `gorm:"not null" json:"nickname"`
    Email             string         `gorm:"unique;not null" json:"email"`
    PasswordHash      string         `gorm:"not null" json:"-"`
    CreatedAt         time.Time      `gorm:"default:current_timestamp" json:"created_at"`
    SubscriptionStatus *string       `json:"subscription_status,omitempty"`
}

// BeforeCreate — хук GORM, который генерирует UUID перед созданием записи.
func(u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.ID = uuid.New()
    return
}