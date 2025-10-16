package service

import (
	"context"
	"valibibe/internal/controller/dto"
	apperrors "valibibe/internal/errors"

	"valibibe/internal/repository/interfaces"

	"github.com/google/uuid"
)

type NoteTagService struct {
	noteRepo interfaces.NoteRepository
	tagRepo  interfaces.TagRepository
}

func NewNoteTagService(noteRepo interfaces.NoteRepository, tagRepo interfaces.TagRepository) *NoteTagService {
	return &NoteTagService{noteRepo: noteRepo, tagRepo: tagRepo}
}

// Добавить тег к заметке
func (s *NoteTagService) AddTag(ctx context.Context, userID, noteID, tagID string) error {
	// 1. Проверить, что заметка принадлежит пользователю
	note, err := s.noteRepo.GetNoteByIDAndUserID(ctx, noteID, userID)
	if err != nil {
		return err // DB error
	}
	if note == nil {
		return apperrors.ErrNotFound // Note not found or access denied
	}

	// 2. Проверить, что тег принадлежит пользователю
	uid, _ := uuid.Parse(userID)
	tid, _ := uuid.Parse(tagID)
	tag, err := s.tagRepo.GetByID(ctx, uid, tid)
	if err != nil {
		return err // DB error
	}
	if tag == nil {
		return apperrors.ErrNotFound // Tag not found or access denied
	}

	// 3. Добавить связь
	nid, _ := uuid.Parse(noteID)
	return s.noteRepo.AddTag(ctx, nid, tid)
}

// Удалить тег у заметки
func (s *NoteTagService) RemoveTag(ctx context.Context, userID, noteID, tagID string) error {
	// 1. Проверить, что заметка принадлежит пользователю
	note, err := s.noteRepo.GetNoteByIDAndUserID(ctx, noteID, userID)
	if err != nil {
		return err // DB error
	}
	if note == nil {
		return apperrors.ErrNotFound // Note not found or access denied
	}

	// 2. Проверить, что тег принадлежит пользователю (не обязательно, но для консистентности)
	uid, _ := uuid.Parse(userID)
	tid, _ := uuid.Parse(tagID)
	tag, err := s.tagRepo.GetByID(ctx, uid, tid)
	if err != nil {
		return err // DB error
	}
	if tag == nil {
		return apperrors.ErrNotFound // Tag not found or access denied
	}

	// 3. Удалить связь
	nid, _ := uuid.Parse(noteID)
	return s.noteRepo.RemoveTag(ctx, nid, tid)
}

// Массовое добавление тегов к заметкам (batch upsert)
func (s *NoteTagService) AddTagsBatch(ctx context.Context, userID string, noteTags []dto.NoteTagInput) error {
	if len(noteTags) == 0 {
		return nil
	}

	// 1. Собрать все уникальные ID заметок и тегов для проверки
	noteIDs := make([]string, 0, len(noteTags))
	tagIDs := make([]string, 0, len(noteTags))
	noteIDSet := make(map[string]struct{})
	tagIDSet := make(map[string]struct{})

	for _, nt := range noteTags {
		if _, exists := noteIDSet[nt.NoteID]; !exists {
			noteIDSet[nt.NoteID] = struct{}{}
			noteIDs = append(noteIDs, nt.NoteID)
		}
		if _, exists := tagIDSet[nt.TagID]; !exists {
			tagIDSet[nt.TagID] = struct{}{}
			tagIDs = append(tagIDs, nt.TagID)
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

	// 3. Проверить, что все теги принадлежат пользователю
	tagCount, err := s.tagRepo.CountTagsByIDsAndUserID(ctx, tagIDs, userID)
	if err != nil {
		return err
	}
	if tagCount != len(tagIDs) {
		return apperrors.ErrNotFound // One or more tags not found or access denied
	}

	// 4. Подготовить данные для пакетной вставки
	parsed := make([]interfaces.NoteTag, 0, len(noteTags))
	for _, nt := range noteTags {
		nid, _ := uuid.Parse(nt.NoteID)
		tid, _ := uuid.Parse(nt.TagID)
		parsed = append(parsed, interfaces.NoteTag{NoteID: nid, TagID: tid})
	}

	return s.noteRepo.AddTagsBatch(ctx, parsed)
}
