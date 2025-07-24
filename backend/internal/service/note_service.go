package service

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/google/uuid"
    "my_app_backend/internal/models"
    "my_app_backend/internal/repository/interfaces"
)

type NoteService struct {
    noteRepo interfaces.NoteRepository
}

func NewNoteService(noteRepo interfaces.NoteRepository) *NoteService {
    return &NoteService{noteRepo: noteRepo}
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
    note, err := s.noteRepo.GetNoteByIDAndUserID(ctx, noteID, userID)
    if err != nil {
        return nil, err
    }
    if note == nil {
        return nil, errors.New("note not found or access denied")
    }
    return note, nil
}

func (s *NoteService) GetAllNotesByUserID(ctx context.Context, filter *models.NoteFilter) (*models.PaginatedNotes, error) {
    return s.noteRepo.GetAllNotesByUserID(ctx, filter)
}

func (s *NoteService) UpdateNote(ctx context.Context, userID, noteID string, input *models.NoteInput) (*models.Note, error) {
    note, err := s.noteRepo.GetNoteByIDAndUserID(ctx, noteID, userID)
    if err != nil {
        return nil, err
    }
    if note == nil {
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
    note, err := s.noteRepo.GetNoteByIDAndUserID(ctx, noteID, userID)
    if err != nil {
        return nil, err
    }
    if note == nil {
        return nil, errors.New("note not found or access denied")
    }

    note.Archived = true

    err = s.noteRepo.UpdateNote(ctx, note)
    if err != nil {
        return nil, err
    }

    return note, nil
}

func (s *NoteService) UnArchiveNote(ctx context.Context, userID, noteID string) (*models.Note, error) {
    note, err := s.noteRepo.GetNoteByIDAndUserID(ctx, noteID, userID)
    if err != nil {
        return nil, err
    }
    if note == nil {
        return nil, errors.New("note not found or access denied")
    }

    note.Archived = false

    err = s.noteRepo.UpdateNote(ctx, note)
    if err != nil {
        return nil, err
    }

    return note, nil
}

func (s *NoteService) DeleteNote(ctx context.Context, userID, noteID string) error {
    note, err := s.noteRepo.GetNoteByIDAndUserID(ctx, noteID, userID)
    if err != nil {
        return err
    }
    if note == nil {
        return errors.New("note not found or access denied")
    }

    return s.noteRepo.DeleteNote(ctx, noteID)
}

func (s *NoteService) UpdateMemoryLevel(ctx context.Context, userID, noteID string, remembered bool) error {
    note, err := s.noteRepo.GetNoteByIDAndUserID(ctx, noteID, userID)
    if err != nil {
        return err
    }
    if note == nil {
        return errors.New("note not found or access denied")
    }

    if !remembered {
        note.MemoryLevel = 0
        note.NextReviewAt = nil
    } else {
        note.MemoryLevel += 20
        if note.MemoryLevel > 100 {
            note.MemoryLevel = 100
        }
        nextReview := calcNextReviewAt(note.MemoryLevel)
        fmt.Printf("Calculated next review: %v\n", nextReview)
        note.NextReviewAt = nextReview
    }

    err = s.noteRepo.UpdateNote(ctx, note)
    if err != nil {
        return err
    }

    updatedNote, _ := s.noteRepo.GetNoteByIDAndUserID(ctx, noteID, userID)
    fmt.Printf("Saved in DB - MemoryLevel: %d, NextReviewAt: %v\n",
        updatedNote.MemoryLevel, updatedNote.NextReviewAt)

    return nil
}

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
