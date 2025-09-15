package repository

import (
    "context"
    "errors"

    "gorm.io/gorm"
    "valibibe/internal/models"
	"valibibe/internal/repository/interfaces"
)

type folderRepo struct {
    db *gorm.DB
}

func NewFolderRepo(db *gorm.DB) interfaces.FolderRepository {
    return &folderRepo{db: db}
}

// Create new folder
func (r *folderRepo) Create(ctx context.Context, folder *models.Folder) error {
    return r.db.WithContext(ctx).Create(folder).Error
}

// Get folder by ID (with user check)
func (r *folderRepo) GetByID(ctx context.Context, userID, id string) (*models.Folder, error) {
    var folder models.Folder

    err := r.db.WithContext(ctx).
           Where("id = ? AND user_id = ?", id, userID).
           First(&folder).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound){
            return nil, nil
        }
        return nil, err
    }

    return &folder, nil
}

// List all folders of user
func (r *folderRepo) ListByUser(ctx context.Context, userID string) ([]models.Folder, error) {
	var folders []models.Folder
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&folders).Error
	if err != nil {
		return nil, err
	}
	return folders, nil
}

// Update folder
func (r *folderRepo) Update(ctx context.Context, folder *models.Folder) error {
	return r.db.WithContext(ctx).Save(folder).Error
}

// Delete folder (will cascade delete children + notes)
func (r *folderRepo) Delete(ctx context.Context, userID, id string) error {
	// доп. защита: удаляем только если принадлежит юзеру
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Folder{}).Error
}

// Check if candidateID is descendant of ancestorID
// Используем рекурсивный CTE (только для PostgreSQL).
func (r *folderRepo) IsDescendant(ctx context.Context, userID, ancestorID, candidateID string) (bool, error) {
	var exists bool
	query := `
		WITH RECURSIVE subfolders AS (
			SELECT id, parent_id
			FROM folders
			WHERE id = ? AND user_id = ?
			UNION
			SELECT f.id, f.parent_id
			FROM folders f
			INNER JOIN subfolders sf ON sf.id = f.parent_id
			WHERE f.user_id = ?
		)
		SELECT EXISTS (
			SELECT 1 FROM subfolders WHERE id = ?
		)
	`
	err := r.db.WithContext(ctx).
		Raw(query, ancestorID, userID, userID, candidateID).
		Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}