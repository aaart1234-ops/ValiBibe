package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"valibibe/internal/controller/dto"
	"valibibe/internal/repository/interfaces"
)

type ReviewSessionService struct {
	noteRepo interfaces.NoteRepository
}

func NewReviewSessionService(noteRepo interfaces.NoteRepository) *ReviewSessionService {
	return &ReviewSessionService{
		noteRepo: noteRepo,
	}
}

// CreateReviewSession создает сессию повторения с фильтрами
func (s *ReviewSessionService) CreateReviewSession(ctx context.Context, userID string, input *dto.ReviewSessionInput) (*dto.ReviewSessionResponse, error) {
	// Валидация входных данных
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	// Конвертируем userID в UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// Получаем заметки для повторения
	notes, err := s.noteRepo.GetNotesForReview(ctx, userUUID, input)
	if err != nil {
		return nil, err
	}

	// Конвертируем заметки в формат ответа
	reviewNotes := make([]dto.ReviewSessionNote, len(notes))
	for i, note := range notes {
		reviewNotes[i] = dto.ReviewSessionNote{
			ID:          note.ID.String(),
			Title:       note.Title,
			Content:     note.Content,
			MemoryLevel: note.MemoryLevel,
			CreatedAt:   note.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   note.UpdatedAt.Format(time.RFC3339),
		}

		// Добавляем next_review_at если есть
		if note.NextReviewAt != nil {
			reviewNotes[i].NextReviewAt = note.NextReviewAt.Format(time.RFC3339)
		}

		// Добавляем информацию о папке
		if note.Folder != nil {
			reviewNotes[i].FolderID = note.Folder.ID.String()
			reviewNotes[i].FolderName = note.Folder.Name
		}

		// Добавляем теги
		reviewNotes[i].Tags = make([]dto.Tag, len(note.Tags))
		for j, tag := range note.Tags {
			reviewNotes[i].Tags[j] = dto.Tag{
				ID:   tag.ID.String(),
				Name: tag.Name,
			}
		}
	}

	return &dto.ReviewSessionResponse{
		Notes: reviewNotes,
		Total: len(reviewNotes),
	}, nil
}
