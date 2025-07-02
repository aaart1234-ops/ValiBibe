package models

type PaginatedNotes struct {
    Notes []Note `json:"notes"`
    Total int64 `json:"total"`
}