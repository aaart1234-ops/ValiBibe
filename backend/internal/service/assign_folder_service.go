package service

import (
	"context"
	apperrors "valibibe/internal/errors"

	"github.com/google/uuid"
	"valibibe/internal/repository/interfaces"
)

type AssignFolderService struct {
	noteRepo   interfaces.NoteRepository
	folderRepo interfaces.FolderRepository
}

func NewAssignFolderService(noteRepo interfaces.NoteRepository, folderRepo interfaces.FolderRepository) *AssignFolderService {
	return &AssignFolderService{
		noteRepo:   noteRepo,
		folderRepo: folderRepo,
	}
}

func (s *AssignFolderService) AssignFolder(ctx context.Context, userID, noteID, folderID string) error {
	// 1. Проверить, что заметка принадлежит пользователю
	note, err := s.noteRepo.GetNoteByIDAndUserID(ctx, noteID, userID)
	if err != nil {
		return err // DB error
	}
	if note == nil {
		return apperrors.ErrNotFound // Note not found or access denied
	}

	// 2. Проверить, что папка принадлежит пользователю
	folder, err := s.folderRepo.GetByID(ctx, userID, folderID)
	if err != nil {
		return err // DB error
	}
	if folder == nil {
		return apperrors.ErrNotFound // Folder not found or access denied
	}

	// 3. Обновить папку у заметки
	fid, _ := uuid.Parse(folderID)
	return s.noteRepo.UpdateFolder(ctx, note.ID, &fid)
}

func (s *AssignFolderService) RemoveFolder(ctx context.Context, userID, noteID, folderID string) error {
	// 1. Проверить, что заметка принадлежит пользователю
	note, err := s.noteRepo.GetNoteByIDAndUserID(ctx, noteID, userID)
	if err != nil {
		return err // DB error
	}
	if note == nil {
		return apperrors.ErrNotFound // Note not found or access denied
	}

	// 2. Проверить, что папка принадлежит пользователю (хотя бы для консистентности)
	folder, err := s.folderRepo.GetByID(ctx, userID, folderID)
	if err != nil {
		return err // DB error
	}
	if folder == nil {
		return apperrors.ErrNotFound // Folder not found or access denied
	}

	// 3. Убрать папку
	return s.noteRepo.UpdateFolder(ctx, note.ID, nil)
}

func (s *AssignFolderService) BatchAssignFolder(ctx context.Context, userID string, noteIDs []string, folderID *string) error {
	// 1. Если папка указана, проверить, что она принадлежит пользователю
	if folderID != nil && *folderID != "" {
		folder, err := s.folderRepo.GetByID(ctx, userID, *folderID)
		if err != nil {
			return err // DB error
		}
		if folder == nil {
			return apperrors.ErrNotFound // Folder not found or access denied
		}
	}

	// 2. Проверить, что все заметки принадлежат пользователю
	count, err := s.noteRepo.CountNotesByIDsAndUserID(ctx, noteIDs, userID)
	if err != nil {
		return err
	}
	if count != len(noteIDs) {
		return apperrors.ErrNotFound // One or more notes not found or access denied
	}

	// 3. Массово обновить
	var fUUID *uuid.UUID
	if folderID != nil && *folderID != "" {
		fid, _ := uuid.Parse(*folderID)
		fUUID = &fid
	}

	return s.noteRepo.BatchUpdateFolder(ctx, noteIDs, fUUID)
}