package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"valibibe/internal/controller/dto"
	"valibibe/internal/models"
	"valibibe/internal/repository/interfaces"
)

type NoteRepo struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) interfaces.NoteRepository {
	return &NoteRepo{db: db}
}

func (r *NoteRepo) CreateNote(ctx context.Context, note *models.Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *NoteRepo) GetNoteByID(ctx context.Context, userID, noteID uuid.UUID) (*models.Note, error) {
	var note models.Note
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", noteID, userID).
		First(&note).Error

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

func (r *NoteRepo) GetAllNotesByUserID(ctx context.Context, filter *dto.NoteFilter) (*dto.PaginatedNotes, error) {
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

	if filter.FolderID != nil && *filter.FolderID != "" {
		query = query.Where("folder_id = ?", *filter.FolderID)
	}

	if len(filter.TagIDs) > 0 {
		query = query.Joins("JOIN note_tags nt ON nt.note_id = notes.id").
			Where("nt.tag_id IN ?", filter.TagIDs)
		// Если нужно только заметки, у которых есть все указанные теги, используйте HAVING COUNT(DISTINCT nt.tag_id) = len(filter.TagIDs)
		// query = query.Group("notes.id").Having("COUNT(DISTINCT nt.tag_id) = ?", len(filter.TagIDs))
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

	// Preload tags to avoid N+1
	query = query.Preload("Tags")

	if err := query.Find(&notes).Error; err != nil {
		return nil, err
	}

	return &dto.PaginatedNotes{
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

func (r *NoteRepo) UpdateFolder(ctx context.Context, userID, noteID uuid.UUID, folderID *uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Note{}).
		Where("id = ? AND user_id = ?", noteID, userID).
		Update("folder_id", folderID).Error
}

func (r *NoteRepo) BatchUpdateFolder(ctx context.Context, userID uuid.UUID, noteIDs []uuid.UUID, folderID *uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Note{}).
		Where("id IN ? AND user_id = ?", noteIDs, userID).
		Update("folder_id", folderID).Error
}

// Добавить тег к заметке (upsert)
func (r *NoteRepo) AddTag(ctx context.Context, noteID, tagID uuid.UUID) error {
	return r.db.WithContext(ctx).Exec(`
        INSERT INTO note_tags (note_id, tag_id)
        VALUES (?, ?)
        ON CONFLICT (note_id, tag_id) DO NOTHING
    `, noteID, tagID).Error
}

// Удалить тег у заметки
func (r *NoteRepo) RemoveTag(ctx context.Context, noteID, tagID uuid.UUID) error {
	return r.db.WithContext(ctx).Exec(`
        DELETE FROM note_tags WHERE note_id = ? AND tag_id = ?
    `, noteID, tagID).Error
}

// Массовая вставка (batch upsert)
func (r *NoteRepo) AddTagsBatch(ctx context.Context, noteTags []interfaces.NoteTag) error {
	if len(noteTags) == 0 {
		return nil
	}

	values := make([]string, 0, len(noteTags))
	args := make([]interface{}, 0, len(noteTags)*2)
	for _, nt := range noteTags {
		values = append(values, "(?, ?  )")
		args = append(args, nt.NoteID, nt.TagID)
	}

	query := fmt.Sprintf(`
        INSERT INTO note_tags (note_id, tag_id)
        VALUES %s
        ON CONFLICT (note_id, tag_id) DO NOTHING
    `, strings.Join(values, ","))

	return r.db.WithContext(ctx).Exec(query, args...).Error
}
