package dto

import (
    "valibibe/internal/models"
)

type PaginatedNotes struct {
    Notes []models.Note `json:"notes"`
    Total int64 `json:"total"`
}