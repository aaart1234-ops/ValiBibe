package interfaces

import (
	"context"

	"valibibe/internal/models"
)

type FolderRepository interface {
	Create(ctx context.Context, folder *models.Folder) error
	GetByID(ctx context.Context, userID, id string) (*models.Folder, error)
	ListByUser(ctx context.Context, userID string) ([]models.Folder, error)
	Update(ctx context.Context, folder *models.Folder) error
	Delete(ctx context.Context, userID, id string) error
	IsDescendant(ctx context.Context, userID, ancestorID, candidateID string) (bool, error)
}
