package repository

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "gorm.io/gorm"

    "valibibe/internal/models"
    "valibibe/internal/repository/interfaces"
)

type tagRepository struct {
    db *gorm.DB
}

func NewTagRepository(db *gorm.DB) interfaces.TagRepository {
    return &tagRepository{db: db}
}

func (r *tagRepository) Create(ctx context.Context, tag *models.Tag) error {
    return r.db.WithContext(ctx).Create(tag).Error
}

func (r *tagRepository) GetByID(ctx context.Context, userID, tagID uuid.UUID) (*models.Tag, error) {
    var tag models.Tag
    err := r.db.WithContext(ctx).
        Where("id = ? AND user_id = ?", tagID, userID).
        First(&tag).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }

    return &tag, err
}

func (r *tagRepository) Update(ctx context.Context, tag *models.Tag) error {
    result := r.db.WithContext(ctx).
        Model(&models.Tag{}).
        Where("id = ? AND user_id = ?", tag.ID, tag.UserID).
        Updates(map[string]interface{}{
            "name": tag.Name,
        })

    if result.RowsAffected == 0 {
        return errors.New("tag not found or not owned by user")
    }

    return result.Error
}

func (r *tagRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Tag, error) {
    var tags []models.Tag
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Order("created_at DESC").
        Find(&tags).Error
    return tags, err
}

func (r *tagRepository) Delete(ctx context.Context, userID, tagID uuid.UUID) error {
    result := r.db.WithContext(ctx).
        Where("id = ? AND user_id = ?", tagID, userID).
        Delete(&models.Tag{})

    if result.RowsAffected == 0 {
        return errors.New("tag not found or not owned by user")
    }

    return result.Error
}

func (r *tagRepository) ExistsByName(ctx context.Context, userID uuid.UUID, name string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&models.Tag{}).
        Where("user_id = ? AND LOWER(name) = LOWER(?)", userID, name).
        Count(&count).Error

    return count > 0, err
}

// ------------------- связи тегов и заметок -------------------

func (r *tagRepository) AttachToNote(ctx context.Context, noteID, tagID uuid.UUID) error {
    relation := models.NoteTag{NoteID: noteID, TagID: tagID}
    return r.db.WithContext(ctx).Create(&relation).Error
}

func (r *tagRepository) DetachFromNote(ctx context.Context, noteID, tagID uuid.UUID) error {
    return r.db.WithContext(ctx).
        Where("note_id = ? AND tag_id = ?", noteID, tagID).
        Delete(&models.NoteTag{}).Error
}

func (r *tagRepository) ListTagsByNote(ctx context.Context, noteID uuid.UUID) ([]models.Tag, error) {
    var tags []models.Tag
    err := r.db.WithContext(ctx).
        Joins("JOIN note_tags nt ON nt.tag_id = tags.id").
        Where("nt.note_id = ?", noteID).
        Find(&tags).Error
    return tags, err
}

func (r *tagRepository) ListNotesByTag(ctx context.Context, userID, tagID uuid.UUID) ([]models.Note, error) {
    var notes []models.Note
    err := r.db.WithContext(ctx).
        Joins("JOIN note_tags nt ON nt.note_id = notes.id").
        Where("nt.tag_id = ? AND notes.user_id = ?", tagID, userID).
        Find(&notes).Error
    return notes, err
}
