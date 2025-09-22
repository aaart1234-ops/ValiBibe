package interfaces

import (
    "context"

    "github.com/google/uuid"
    "valibibe/internal/models"
)

type TagRepository interface {
    Create(ctx context.Context, tag *models.Tag) error
    GetByID(ctx context.Context, userID, tagID uuid.UUID) (*models.Tag, error)
    ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Tag, error)
    Update(ctx context.Context, tag *models.Tag) error
    Delete(ctx context.Context, userID, tagID uuid.UUID) error
    ExistsByName(ctx context.Context, userID uuid.UUID, name string) (bool, error)

    // связи заметок и тегов
    AttachToNote(ctx context.Context, noteID, tagID uuid.UUID) error
    DetachFromNote(ctx context.Context, noteID, tagID uuid.UUID) error
    ListTagsByNote(ctx context.Context, noteID uuid.UUID) ([]models.Tag, error)
    ListNotesByTag(ctx context.Context, userID, tagID uuid.UUID) ([]models.Note, error)
}
