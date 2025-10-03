package unit

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"valibibe/internal/controller/dto"
	"valibibe/internal/models"
	"valibibe/internal/repository/interfaces"
	"valibibe/internal/service"
)

// MockNoteRepo реализует интерфейс NoteRepository для моков
type MockNoteRepo struct {
	mock.Mock
}

func (m *MockNoteRepo) CreateNote(ctx context.Context, note *models.Note) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockNoteRepo) GetNoteByIDAndUserID(ctx context.Context, noteID string, userID string) (*models.Note, error) {
	args := m.Called(ctx, noteID, userID)
	note, _ := args.Get(0).(*models.Note)
	return note, args.Error(1)
}

func (m *MockNoteRepo) GetAllNotesByUserID(ctx context.Context, filter *dto.NoteFilter) (*dto.PaginatedNotes, error) {
	args := m.Called(ctx, filter)

	// Безопасно извлекаем PaginatedNotes
	result, ok := args.Get(0).(*dto.PaginatedNotes)
	if !ok && args.Get(0) != nil {
		return nil, fmt.Errorf("expected *dto.PaginatedNotes, got %T", args.Get(0))
	}

	return result, args.Error(1)
}

func (m *MockNoteRepo) UpdateNote(ctx context.Context, note *models.Note) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockNoteRepo) DeleteNote(ctx context.Context, noteID string) error {
	args := m.Called(ctx, noteID)
	return args.Error(0)
}

func (m *MockNoteRepo) ArchiveNote(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNoteRepo) UnArchiveNote(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNoteRepo) GetNoteByID(ctx context.Context, userID, noteID uuid.UUID) (*models.Note, error) {
	args := m.Called(ctx, userID, noteID)
	note, _ := args.Get(0).(*models.Note)
	return note, args.Error(1)
}

func (m *MockNoteRepo) UpdateFolder(ctx context.Context, userID, noteID uuid.UUID, folderID *uuid.UUID) error {
	args := m.Called(ctx, userID, noteID, folderID)
	return args.Error(0)
}

func (m *MockNoteRepo) BatchUpdateFolder(ctx context.Context, userID uuid.UUID, noteIDs []uuid.UUID, folderID *uuid.UUID) error {
	args := m.Called(ctx, userID, noteIDs, folderID)
	return args.Error(0)
}

func (m *MockNoteRepo) AddTag(ctx context.Context, noteID, tagID uuid.UUID) error {
	args := m.Called(ctx, noteID, tagID)
	return args.Error(0)
}

func (m *MockNoteRepo) RemoveTag(ctx context.Context, noteID, tagID uuid.UUID) error {
	args := m.Called(ctx, noteID, tagID)
	return args.Error(0)
}

func (m *MockNoteRepo) AddTagsBatch(ctx context.Context, noteTags []interfaces.NoteTag) error {
	args := m.Called(ctx, noteTags)
	return args.Error(0)
}

// ====== Тесты NoteService ======

func TestNoteService_CreateNote(t *testing.T) {
	mockRepo := new(MockNoteRepo)
	noteService := service.NewNoteService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	input := &dto.NoteInput{
		Title:   "Test Title",
		Content: "Test Content",
	}

	mockRepo.On("CreateNote", ctx, mock.MatchedBy(func(note *models.Note) bool {
		return note.Title == input.Title && note.Content == input.Content && note.UserID == userID
	})).Return(nil)

	createdNote, err := noteService.CreateNote(ctx, userID.String(), input)
	assert.NoError(t, err)
	assert.Equal(t, input.Title, createdNote.Title)
	assert.Equal(t, input.Content, createdNote.Content)
	assert.Equal(t, userID, createdNote.UserID)

	mockRepo.AssertExpectations(t)
}

func TestNoteService_GetNoteByID(t *testing.T) {
	mockRepo := new(MockNoteRepo)
	noteService := service.NewNoteService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	noteID := uuid.New()

	expectedNote := &models.Note{
		ID:      noteID,
		Title:   "Test Title",
		Content: "Test Content",
		UserID:  userID,
	}

	//mockRepo.On("GetNoteByID", ctx, noteID.String()).Return(expectedNote, nil)
	mockRepo.On("GetNoteByIDAndUserID", ctx, noteID.String(), userID.String()).Return(expectedNote, nil)
	note, err := noteService.GetNoteByID(ctx, userID.String(), noteID.String())
	assert.NoError(t, err)
	assert.Equal(t, expectedNote, note)

	mockRepo.AssertExpectations(t)
}

func TestNoteService_GetAllNotesByUserID(t *testing.T) {
	mockRepo := new(MockNoteRepo)
	noteService := service.NewNoteService(mockRepo)
	ctx := context.Background()

	userID := uuid.New().String()
	filter := &dto.NoteFilter{
		UserID: userID,
		Search: "Note",
		SortBy: "created_at",
		Order:  "desc",
	}

	notes := []models.Note{
		{ID: uuid.New(), Title: "Note 1", Content: "Content 1", UserID: uuid.MustParse(userID)},
		{ID: uuid.New(), Title: "Note 2", Content: "Content 2", UserID: uuid.MustParse(userID)},
	}

	mockRepo.On("GetAllNotesByUserID", ctx, filter).Return(&dto.PaginatedNotes{
		Notes: notes,
		Total: int64(len(notes)),
	}, nil)

	result, err := noteService.GetAllNotesByUserID(ctx, filter)
	assert.NoError(t, err)
	assert.Equal(t, notes, result.Notes)
	assert.Equal(t, int64(2), result.Total)

	mockRepo.AssertExpectations(t)
}

func TestNoteService_UpdateNote(t *testing.T) {
	mockRepo := new(MockNoteRepo)
	noteService := service.NewNoteService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	noteID := uuid.New()

	note := &models.Note{
		ID:      noteID,
		Title:   "Old Title",
		Content: "Old Content",
		UserID:  userID,
	}

	input := &dto.NoteInput{
		Title:   "Updated Title",
		Content: "Updated Content",
	}

	//mockRepo.On("GetNoteByID", ctx, noteID.String()).Return(note, nil)
	mockRepo.On("GetNoteByIDAndUserID", ctx, noteID.String(), userID.String()).Return(note, nil)
	mockRepo.On("UpdateNote", ctx, mock.MatchedBy(func(n *models.Note) bool {
		return n.Title == input.Title && n.Content == input.Content && n.ID == noteID
	})).Return(nil)

	updatedNote, err := noteService.UpdateNote(ctx, userID.String(), noteID.String(), input)
	assert.NoError(t, err)
	assert.Equal(t, input.Title, updatedNote.Title)
	assert.Equal(t, input.Content, updatedNote.Content)

	mockRepo.AssertExpectations(t)
}

func TestNoteService_DeleteNote(t *testing.T) {
	mockRepo := new(MockNoteRepo)
	noteService := service.NewNoteService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	noteID := uuid.New()

	note := &models.Note{
		ID:     noteID,
		UserID: userID,
	}

	//mockRepo.On("GetNoteByID", ctx, noteID.String()).Return(note, nil)
	mockRepo.On("GetNoteByIDAndUserID", ctx, noteID.String(), userID.String()).Return(note, nil)
	mockRepo.On("DeleteNote", ctx, noteID.String()).Return(nil)

	err := noteService.DeleteNote(ctx, userID.String(), noteID.String())
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestNoteService_ArchiveNote(t *testing.T) {
	mockRepo := new(MockNoteRepo)
	noteService := service.NewNoteService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	noteID := uuid.New()

	// Исходная заметка (до архивации)
	note := &models.Note{
		ID:       noteID,
		UserID:   userID,
		Archived: false, // Изначально не архивирована
	}

	// Ожидаемая заметка после архивации
	archivedNote := &models.Note{
		ID:       noteID,
		UserID:   userID,
		Archived: true, // Теперь архивирована
	}

	// Мокируем вызовы репозитория
	//mockRepo.On("GetNoteByID", ctx, noteID.String()).Return(note, nil)
	mockRepo.On("GetNoteByIDAndUserID", ctx, noteID.String(), userID.String()).Return(note, nil)
	mockRepo.On("UpdateNote", ctx, archivedNote).Return(nil)

	// Вызываем метод и проверяем, что ошибки нет и заметка вернулась с Archived = true
	updatedNote, err := noteService.ArchiveNote(ctx, userID.String(), noteID.String())

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, updatedNote)
	assert.True(t, updatedNote.Archived)    // Убеждаемся, что заметка архивирована
	assert.Equal(t, noteID, updatedNote.ID) // Проверяем, что ID совпадает

	// Проверяем, что все ожидаемые вызовы репозитория были выполнены
	mockRepo.AssertExpectations(t)
}

func TestNoteService_UpdateMemoryLevel(t *testing.T) {
	mockRepo := new(MockNoteRepo)
	noteService := service.NewNoteService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	noteID := uuid.New()

	note := &models.Note{
		ID:          noteID,
		UserID:      userID,
		MemoryLevel: 40,
	}

	//mockRepo.On("GetNoteByID", ctx, noteID.String()).Return(note, nil)
	mockRepo.On("GetNoteByIDAndUserID", ctx, noteID.String(), userID.String()).Return(note, nil)
	mockRepo.On("UpdateNote", ctx, mock.MatchedBy(func(n *models.Note) bool {
		return n.ID == noteID && (n.MemoryLevel == 60 || n.MemoryLevel == 0)
	})).Return(nil)

	// Тестируем рост memoryLevel
	err := noteService.UpdateMemoryLevel(ctx, userID.String(), noteID.String(), true)
	assert.NoError(t, err)
	assert.Equal(t, 60, note.MemoryLevel)
	assert.NotNil(t, note.NextReviewAt)

	// Тестируем сброс memoryLevel
	err = noteService.UpdateMemoryLevel(ctx, userID.String(), noteID.String(), false)
	assert.NoError(t, err)
	assert.Equal(t, 0, note.MemoryLevel)
	assert.Nil(t, note.NextReviewAt)

	mockRepo.AssertExpectations(t)
}
