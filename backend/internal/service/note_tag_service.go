package service

import (
	"context"
	"errors"

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
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid userID")
	}
	nid, err := uuid.Parse(noteID)
	if err != nil {
		return errors.New("invalid noteID")
	}
	tid, err := uuid.Parse(tagID)
	if err != nil {
		return errors.New("invalid tagID")
	}

	// Проверим, что заметка принадлежит пользователю
	note, err := s.noteRepo.GetNoteByID(ctx, uid, nid)
	if err != nil {
		return err
	}
	if note == nil {
		return errors.New("note not found")
	}

	// Проверим, что тег принадлежит пользователю
	tag, err := s.tagRepo.GetByID(ctx, uid, tid)
	if err != nil {
		return err
	}
	if tag == nil {
		return errors.New("tag not found")
	}

	return s.noteRepo.AddTag(ctx, nid, tid)
}

// Удалить тег у заметки
func (s *NoteTagService) RemoveTag(ctx context.Context, userID, noteID, tagID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid userID")
	}
	nid, err := uuid.Parse(noteID)
	if err != nil {
		return errors.New("invalid noteID")
	}
	tid, err := uuid.Parse(tagID)
	if err != nil {
		return errors.New("invalid tagID")
	}

	// Проверим, что заметка принадлежит пользователю
	note, err := s.noteRepo.GetNoteByID(ctx, uid, nid)
	if err != nil {
		return err
	}
	if note == nil {
		return errors.New("note not found")
	}

	return s.noteRepo.RemoveTag(ctx, nid, tid)
}

// Массовое добавление тегов к заметкам (batch upsert)
func (s *NoteTagService) AddTagsBatch(ctx context.Context, userID string, noteTags []struct {
	NoteID string
	TagID  string
}) error {
	/*uid, err := uuid.Parse(userID)
	  if err != nil {
	      return errors.New("invalid userID")
	  }*/

	parsed := make([]interfaces.NoteTag, 0, len(noteTags))

	for _, nt := range noteTags {
		nid, err := uuid.Parse(nt.NoteID)
		if err != nil {
			return errors.New("invalid noteID in batch")
		}
		tid, err := uuid.Parse(nt.TagID)
		if err != nil {
			return errors.New("invalid tagID in batch")
		}

		// (опционально) можно проверять что note и tag принадлежат userID
		parsed = append(parsed, interfaces.NoteTag{NoteID: nid, TagID: tid})
	}

	return s.noteRepo.AddTagsBatch(ctx, parsed)
}
