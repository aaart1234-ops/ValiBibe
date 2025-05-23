package service

import (
    "fmt"
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
	"my_app_backend/internal/models"
    "my_app_backend/internal/repository/interfaces"
)

type NoteService struct {
    noteRepo interfaces.NoteRepository
}

func NewNoteService(noteRepo interfaces.NoteRepository) *NoteService{
    return &NoteService{noteRepo: noteRepo}
}

func (s *NoteService) UpdateMemoryLevel(ctx context.Context, userID, noteID string, remembered bool) error {
    note, err := s.noteRepo.GetNoteByID(ctx, noteID)
    if err != nil {
        return err
    }
    if note == nil || note.UserID.String() != userID {
        return errors.New("note not found or access denied")
    }

    if !remembered {
        // Пользователь не помнит — сброс
        note.MemoryLevel = 0
        note.NextReviewAt = nil
    } else {
        // Пользователь помнит — рост уровня
        note.MemoryLevel += 20
        if note.MemoryLevel > 100 {
            note.MemoryLevel = 100
        }
        nextReview := calcNextReviewAt(note.MemoryLevel)
        fmt.Printf("Calculated next review: %v\n", nextReview) // Логирование
        note.NextReviewAt = nextReview
    }

    err = s.noteRepo.UpdateNote(ctx, note)
    if err != nil {
        return err
    }

    // Проверяем, что сохранилось в БД
    updatedNote, _ := s.noteRepo.GetNoteByID(ctx, noteID)
    fmt.Printf("Saved in DB - MemoryLevel: %d, NextReviewAt: %v\n",
        updatedNote.MemoryLevel, updatedNote.NextReviewAt)

    return nil
}

// Примерная логика: чем выше уровень, тем дольше до следующего повторения
func calcNextReviewAt(level int) *time.Time {
    var days int
    switch {
        case level < 20:
            days = 1
        case level < 40:
            days = 3
        case level < 60:
            days = 5
        case level < 80:
            days = 10
        default:
            days = 30
    }
    t := time.Now().AddDate(0, 0, days)
    return &t
}

func (s *NoteService) CreateNote(ctx context.Context, userID string, input *models.NoteInput) (*models.Note, error) {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return nil, err
    }

    note := &models.Note{
        UserID:      uid,
        Title:       input.Title,
        Content:     input.Content,
        MemoryLevel: 0,
        Archived:    false,
    }

    err = s.noteRepo.CreateNote(ctx, note)
    if err != nil {
        return nil, err
    }

    return note, nil
}

func (s *NoteService) GetNoteByID(ctx context.Context, userID, noteID string) (*models.Note, error) {
    note, err := s.noteRepo.GetNoteByID(ctx, noteID)
    if err != nil {
        return nil, err
    }

    if note == nil || note.UserID.String() != userID {
        return nil, errors.New("note not found or access denied")
    }

    return note, nil
}

func (s *NoteService) GetAllNotesByUserID(ctx context.Context, userID string) ([]models.Note, error) {
    return s.noteRepo.GetAllNotesByUserID(ctx, userID)
}

func (s *NoteService) UpdateNote(ctx context.Context, userID, noteID string, input *models.NoteInput) (*models.Note, error) {
    note, err := s.noteRepo.GetNoteByID(ctx, noteID)
    if err != nil {
        return nil, err
    }

    if note == nil || note.UserID.String() != userID {
        return nil, errors.New("note not found or access denied")
    }

    note.Title = input.Title
    note.Content = input.Content

    err = s.noteRepo.UpdateNote(ctx, note)
    if err != nil {
        return nil, err
    }

    return note, nil
}

func (s *NoteService) ArchiveNote(ctx context.Context, userID, noteID string) (*models.Note, error) {
    note, err := s.noteRepo.GetNoteByID(ctx, noteID)
    if err != nil {
            return nil, err
    }

    if note == nil || note.UserID.String() != userID {
        return nil, errors.New("note not found or access denied")
    }

    note.Archived = true

    err = s.noteRepo.UpdateNote(ctx, note)
    if err != nil {
        return nil, err
    }

    return note, nil  // Возвращаем обновлённую заметку
}

func (s *NoteService) DeleteNote(ctx context.Context, userID, noteID string) error {
    note, err := s.noteRepo.GetNoteByID(ctx, noteID)
    if err != nil {
        return err
    }

    if note == nil || note.UserID.String() != userID {
        return errors.New("note not found or access denied")
    }

    return s.noteRepo.DeleteNote(ctx, noteID)
}