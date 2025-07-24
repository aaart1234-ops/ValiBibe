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
	err := r.db.WithContext(ctx).
		First(&note, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &note, err
}

func (r *NoteRepo) GetNoteByIDAndUserID(ctx context.Context, id string, userID string) (*models.Note, error) {
	var note models.Note
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&note).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &note, err
}

func (r *NoteRepo) GetAllNotesByUserID(ctx context.Context, filter *models.NoteFilter) (*models.PaginatedNotes, error) {
	var (
		notes []models.Note
		total int64
	)

	query := r.db.WithContext(ctx).
		Model(&models.Note{}).
		Where("user_id = ?", filter.UserID)

	if filter.Archived != nil {
	    query = query.Where("archived = ?", *filter.Archived)
	}

	if filter.Search != "" {
		query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(filter.Search)+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	sortField := map[string]string{
		"created_at":     "created_at",
		"next_review_at": "next_review_at",
	}[filter.SortBy]
	if sortField == "" {
		sortField = "created_at"
	}

	order := "desc"
	if strings.ToLower(filter.Order) == "asc" {
		order = "asc"
	}

	query = query.Order(sortField + " " + order)

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset >= 0 {
		query = query.Offset(filter.Offset)
	}

	if err := query.Find(&notes).Error; err != nil {
		return nil, err
	}

	return &models.PaginatedNotes{
		Notes: notes,
		Total: total,
	}, nil
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

func (r *NoteRepo) UnArchiveNote(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&models.Note{}).
		Where("id = ?", id).
		Update("archived", false).Error
}

func (r *NoteRepo) DeleteNote(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Note{}).Error
}
