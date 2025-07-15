package interfaces

import (
	"context"

	"my_app_backend/internal/models"
)

type NoteRepository interface {
    CreateNote(ctx context.Context, note *models.Note) error
    GetNoteByID(ctx context.Context, id string) (*models.Note, error)
    GetNoteByIDAndUserID(ctx context.Context, noteID string, userID string) (*models.Note, error)
    GetAllNotesByUserID(ctx context.Context, filter *models.NoteFilter) (*models.PaginatedNotes, error)
    UpdateNote(ctx context.Context, note *models.Note) error
    ArchiveNote(ctx context.Context, id string) error
    DeleteNote(ctx context.Context, id string) error
}