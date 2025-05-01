package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"

    "my_app_backend/internal/utils"
)

type Note struct {
    ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
    UserID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
    Title        string     `gorm:"type:varchar(255);not null" json:"title"`
    Content      string     `gorm:"type:text" json:"content"`
    MemoryLevel  int        `gorm:"type:int;default:0;check:memory_level >= 0 AND memory_level <= 100" json:"memory_level"`
    Archived     bool       `gorm:"default:false" json:"archived"`
    NextReviewAt *time.Time `json:"next_review_at, omitempty"`
    CreatedAt    *time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAy    *time.Time `gorm:"autoCreateTime" json:"updated_at"`
}

func (n *Note) BeforeCreate(tx *gorm.DB) (err error) {
    return utils.SetUUIDIfNil(&n.ID)(tx)
}