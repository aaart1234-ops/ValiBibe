package interfaces

import (
	"context"

    "github.com/google/uuid"

	"valibibe/internal/models"
	"valibibe/internal/controller/dto"
)

// NoteTag представляет связь между заметкой и тегом для пакетных операций
type NoteTag struct {
    NoteID uuid.UUID
    TagID  uuid.UUID
}

type NoteRepository interface {
    CreateNote(ctx context.Context, note *models.Note) error
    GetNoteByID(ctx context.Context, userID, noteID uuid.UUID) (*models.Note, error)
    GetNoteByIDAndUserID(ctx context.Context, noteID string, userID string) (*models.Note, error)
    CountNotesByIDsAndUserID(ctx context.Context, noteIDs []string, userID string) (int, error)
    GetAllNotesByUserID(ctx context.Context, filter *dto.NoteFilter) (*dto.PaginatedNotes, error)
    UpdateNote(ctx context.Context, note *models.Note) error
    ArchiveNote(ctx context.Context, id string) error
    UnArchiveNote(ctx context.Context, id string) error
    DeleteNote(ctx context.Context, id string) error
    UpdateFolder(ctx context.Context, noteID uuid.UUID, folderID *uuid.UUID) error
    BatchUpdateFolder(ctx context.Context, noteIDs []string, folderID *uuid.UUID) error
    AddTag(ctx context.Context, noteID, tagID uuid.UUID) error
    RemoveTag(ctx context.Context,noteID, tagID uuid.UUID) error
    AddTagsBatch(ctx context.Context, noteTags []NoteTag) error
    GetNotesForReview(ctx context.Context, userID uuid.UUID, filter *dto.ReviewSessionInput) ([]models.Note, error)
}