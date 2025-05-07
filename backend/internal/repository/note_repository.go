package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"my_app_backend/internal/models"
)

type NoteRepo struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) *NoteRepo {
	return &NoteRepo{db: db}
}

func (r *NoteRepo) CreateNote(ctx context.Context, note *models.Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *NoteRepo) GetNoteByID(ctx context.Context, id string) (*models.Note, error) {
	var note models.Note
	err := r.db.WithContext(ctx).First(&note, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &note, err
}

func (r *NoteRepo) GetAllNotesByUserID(ctx context.Context, userID string) ([]models.Note, error) {
	var notes []models.Note
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND archived = false", userID).
		Find(&notes).Error
	return notes, err
}

func (r *NoteRepo) UpdateNote(ctx context.Context, note *models.Note) error {
	return r.db.WithContext(ctx).Save(note).Error
}

func (r *NoteRepo) ArchiveNote(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&models.Note{}).
		Where("id = ?", id).
		Update("archived", true).Error
}

func (r *NoteRepo) DeleteNote(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Note{}).Error
}
