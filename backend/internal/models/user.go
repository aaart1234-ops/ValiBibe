package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"

    "my_app_backend/internal/utils"
)

// User представляет модель пользователя в базе данных.
type User struct {
    ID                uuid.UUID      `gorm:"type:text;primaryKey" json:"id"`
    Nickname          string         `gorm:"not null" json:"nickname"`
    Email             string         `gorm:"unique;not null" json:"email"`
    PasswordHash      string         `gorm:"not null" json:"-"`
    CreatedAt         time.Time      `gorm:"default:current_timestamp" json:"created_at"`
    SubscriptionStatus *string       `json:"subscription_status,omitempty"`
}

type RegisterRequest struct {
    Email    string `json:"email" example:"test@example.com"`
    Password string `json:"password" example:"123456"`
    Nickname string `json:"nickname" example:"CoolUser"`
}

// LoginRequest - модель запроса для входа
type LoginRequest struct {
    Email    string `json:"email" example:"test@example.com"`
    Password string `json:"password" example:"123456"`
}

// BeforeCreate — хук GORM, который вызывается перед созданием записи
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    return utils.SetUUIDIfNil(&u.ID)(tx)
}
