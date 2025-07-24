package models

type NoteFilter struct {
    UserID string
    Search string
    SortBy string
    Order string
    Limit int
    Offset int
    Archived *bool
}

