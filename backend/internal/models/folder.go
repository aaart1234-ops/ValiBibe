package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "valibibe/internal/utils"
)

type Folder struct {
    ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
    UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"-"`
    Name      string     `gorm:"type:varchar(200);not null" json:"name"`
    ParentID  *uuid.UUID `gorm:"type:uuid;index" json:"parent_id,omitempty"`
    CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (f *Folder) BeforeCreate(tx *gorm.DB) (err error) {
    return utils.SetUUIDIfNil(&f.ID)(tx)
}


