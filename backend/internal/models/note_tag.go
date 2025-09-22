package models

import (
    "github.com/google/uuid"
)


type NoteTag struct {
    NoteID uuid.UUID `gorm:"type:uuid;primaryKey"`
    TagID  uuid.UUID `gorm:"type:uuid;primaryKey"`
}
