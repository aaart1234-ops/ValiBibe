package service

import (
	"context"
	"errors"
	apperrors "valibibe/internal/errors"

	"valibibe/internal/controller/dto"
	"valibibe/internal/models"
	"valibibe/internal/repository/interfaces"

	"github.com/google/uuid"
)

type FolderService struct {
	repo interfaces.FolderRepository
}

func NewFolderService(repo interfaces.FolderRepository) *FolderService {
	return &FolderService{repo: repo}
}

// Create new folder
func (s *FolderService) CreateFolder(ctx context.Context, userID string, input dto.FolderCreateInput) (*models.Folder, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid userID")
	}

	var parentID *uuid.UUID
	if input.ParentID != nil {
		pid, err := uuid.Parse(*input.ParentID)
		if err != nil {
			return nil, errors.New("invalid parentID")
		}

		// Security check: убедиться, что родительская папка принадлежит пользователю
		parentFolder, err := s.repo.GetByID(ctx, userID, *input.ParentID)
		if err != nil {
			return nil, err // DB error
		}
		if parentFolder == nil {
			return nil, apperrors.ErrNotFound // Parent folder not found or access denied
		}
		parentID = &pid
	}

	folder := &models.Folder{
		UserID:   uid,
		Name:     input.Name,
		ParentID: parentID,
	}

	if err := s.repo.Create(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}

// Get folder tree
func (s *FolderService) GetFolderTree(ctx context.Context, userID string) ([]dto.FolderNode, error) {
	folders, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return buildFolderTree(folders), nil
}

// Update folder
func (s *FolderService) UpdateFolder(ctx context.Context, userID, folderID string, input dto.FolderUpdateInput) (*models.Folder, error) {
	folder, err := s.repo.GetByID(ctx, userID, folderID)
	if err != nil {
		return nil, err
	}
	if folder == nil {
		return nil, apperrors.ErrNotFound
	}

	if input.Name != nil {
		folder.Name = *input.Name
	}

	if input.ParentID != nil {
		if *input.ParentID == "" { // Unset parent
			folder.ParentID = nil
		} else {
			pid, err := uuid.Parse(*input.ParentID)
			if err != nil {
				return nil, errors.New("invalid parentID")
			}
			// Security check: убедиться, что новая родительская папка принадлежит пользователю
			parentFolder, err := s.repo.GetByID(ctx, userID, *input.ParentID)
			if err != nil {
				return nil, err // DB error
			}
			if parentFolder == nil {
				return nil, apperrors.ErrNotFound // Parent folder not found or access denied
			}
			folder.ParentID = &pid
		}
	}

	if err := s.repo.Update(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}

// Delete folder (cascade notes + children via DB constraints)
func (s *FolderService) DeleteFolder(ctx context.Context, userID, folderID string) error {
	folder, err := s.repo.GetByID(ctx, userID, folderID)
	if err != nil {
		return err
	}
	if folder == nil {
		return apperrors.ErrNotFound
	}
	return s.repo.Delete(ctx, userID, folderID)
}

// helper: строим дерево из списка
func buildFolderTree(folders []models.Folder) []dto.FolderNode {
	idToNode := make(map[string]*dto.FolderNode)
	var roots []*dto.FolderNode

	// создаём узлы
	for _, f := range folders {
		node := &dto.FolderNode{
			ID:       f.ID.String(),
			Name:     f.Name,
			Children: []*dto.FolderNode{},
		}
		if f.ParentID != nil {
			pid := f.ParentID.String()
			node.ParentID = &pid
		}
		idToNode[f.ID.String()] = node
	}

	// связываем parent-child
	for _, node := range idToNode {
		if node.ParentID != nil {
			if parent, ok := idToNode[*node.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		} else {
			roots = append(roots, node)
		}
	}

	// если надо вернуть []dto.FolderNode (а не []*dto.FolderNode)
	result := make([]dto.FolderNode, 0, len(roots))
	for _, r := range roots {
		result = append(result, *r)
	}

	return result
}
