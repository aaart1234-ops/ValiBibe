package repository

import (
	"context"
	"errors"
	"strings"

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

func (r *NoteRepo) GetAllNotesByUserID(ctx context.Context, filter *models.NoteFilter) ([]models.Note, error) {
	var notes []models.Note

	query := r.db.WithContext(ctx).
		Where("user_id = ? AND archived = false", filter.UserID)

	if filter.Search != "" {
		query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(filter.Search)+"%")
	}

	// Безопасно ограничиваем только разрешённые поля
	sortBy := map[string]string{
		"created_at":      "created_at",
		"next_review_at":  "next_review_at",
	}[filter.SortBy]
	if sortBy == "" {
		sortBy = "created_at"
	}

	order := "desc"
	if strings.ToLower(filter.Order) == "asc" {
		order = "asc"
	}

	err := query.Order(sortBy + " " + order).Find(&notes).Error
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
