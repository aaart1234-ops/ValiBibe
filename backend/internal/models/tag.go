package models

import (
	"time"

	"valibibe/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// BeforeCreate — хук GORM, который вызывается перед созданием записи
func (t *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	return utils.SetUUIDIfNil(&t.ID)(tx)
}
