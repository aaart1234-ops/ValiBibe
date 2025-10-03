package unit

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"valibibe/internal/models"
	"valibibe/internal/repository/interfaces"
	"valibibe/internal/service"
)

// MockTagRepo реализует интерфейс TagRepository для моков
type MockTagRepo struct {
	mock.Mock
}

func (m *MockTagRepo) Create(ctx context.Context, tag *models.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

func (m *MockTagRepo) GetByID(ctx context.Context, userID, tagID uuid.UUID) (*models.Tag, error) {
	args := m.Called(ctx, userID, tagID)
	tag, _ := args.Get(0).(*models.Tag)
	return tag, args.Error(1)
}

func (m *MockTagRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Tag, error) {
	args := m.Called(ctx, userID)
	tags, _ := args.Get(0).([]models.Tag)
	return tags, args.Error(1)
}

func (m *MockTagRepo) Update(ctx context.Context, tag *models.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

func (m *MockTagRepo) Delete(ctx context.Context, userID, tagID uuid.UUID) error {
	args := m.Called(ctx, userID, tagID)
	return args.Error(0)
}

func (m *MockTagRepo) ExistsByName(ctx context.Context, userID uuid.UUID, name string) (bool, error) {
	args := m.Called(ctx, userID, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockTagRepo) AttachToNote(ctx context.Context, noteID, tagID uuid.UUID) error {
	args := m.Called(ctx, noteID, tagID)
	return args.Error(0)
}

func (m *MockTagRepo) DetachFromNote(ctx context.Context, noteID, tagID uuid.UUID) error {
	args := m.Called(ctx, noteID, tagID)
	return args.Error(0)
}

func (m *MockTagRepo) ListTagsByNote(ctx context.Context, noteID uuid.UUID) ([]models.Tag, error) {
	args := m.Called(ctx, noteID)
	tags, _ := args.Get(0).([]models.Tag)
	return tags, args.Error(1)
}

func (m *MockTagRepo) ListNotesByTag(ctx context.Context, userID, tagID uuid.UUID) ([]models.Note, error) {
	args := m.Called(ctx, userID, tagID)
	notes, _ := args.Get(0).([]models.Note)
	return notes, args.Error(1)
}

// ====== Тесты NoteTagService ======

func TestNoteTagService_AddTag(t *testing.T) {
	mockNoteRepo := new(MockNoteRepo)
	mockTagRepo := new(MockTagRepo)
	noteTagService := service.NewNoteTagService(mockNoteRepo, mockTagRepo)
	ctx := context.Background()

	userID := uuid.New()
	noteID := uuid.New()
	tagID := uuid.New()

	note := &models.Note{
		ID:     noteID,
		UserID: userID,
		Title:  "Test Note",
	}

	tag := &models.Tag{
		ID:     tagID,
		UserID: userID,
		Name:   "Test Tag",
	}

	// Мокируем вызовы репозиториев
	mockNoteRepo.On("GetNoteByID", ctx, userID, noteID).Return(note, nil)
	mockTagRepo.On("GetByID", ctx, userID, tagID).Return(tag, nil)
	mockNoteRepo.On("AddTag", ctx, noteID, tagID).Return(nil)

	err := noteTagService.AddTag(ctx, userID.String(), noteID.String(), tagID.String())
	assert.NoError(t, err)

	mockNoteRepo.AssertExpectations(t)
	mockTagRepo.AssertExpectations(t)
}

func TestNoteTagService_AddTag_NoteNotFound(t *testing.T) {
	mockNoteRepo := new(MockNoteRepo)
	mockTagRepo := new(MockTagRepo)
	noteTagService := service.NewNoteTagService(mockNoteRepo, mockTagRepo)
	ctx := context.Background()

	userID := uuid.New()
	noteID := uuid.New()
	tagID := uuid.New()

	// Мокируем что заметка не найдена
	mockNoteRepo.On("GetNoteByID", ctx, userID, noteID).Return(nil, nil)

	err := noteTagService.AddTag(ctx, userID.String(), noteID.String(), tagID.String())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "note not found")

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteTagService_AddTag_TagNotFound(t *testing.T) {
	mockNoteRepo := new(MockNoteRepo)
	mockTagRepo := new(MockTagRepo)
	noteTagService := service.NewNoteTagService(mockNoteRepo, mockTagRepo)
	ctx := context.Background()

	userID := uuid.New()
	noteID := uuid.New()
	tagID := uuid.New()

	note := &models.Note{
		ID:     noteID,
		UserID: userID,
		Title:  "Test Note",
	}

	// Мокируем что заметка найдена, но тег нет
	mockNoteRepo.On("GetNoteByID", ctx, userID, noteID).Return(note, nil)
	mockTagRepo.On("GetByID", ctx, userID, tagID).Return(nil, nil)

	err := noteTagService.AddTag(ctx, userID.String(), noteID.String(), tagID.String())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tag not found")

	mockNoteRepo.AssertExpectations(t)
	mockTagRepo.AssertExpectations(t)
}

func TestNoteTagService_RemoveTag(t *testing.T) {
	mockNoteRepo := new(MockNoteRepo)
	mockTagRepo := new(MockTagRepo)
	noteTagService := service.NewNoteTagService(mockNoteRepo, mockTagRepo)
	ctx := context.Background()

	userID := uuid.New()
	noteID := uuid.New()
	tagID := uuid.New()

	note := &models.Note{
		ID:     noteID,
		UserID: userID,
		Title:  "Test Note",
	}

	// Мокируем вызовы репозиториев
	mockNoteRepo.On("GetNoteByID", ctx, userID, noteID).Return(note, nil)
	mockNoteRepo.On("RemoveTag", ctx, noteID, tagID).Return(nil)

	err := noteTagService.RemoveTag(ctx, userID.String(), noteID.String(), tagID.String())
	assert.NoError(t, err)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteTagService_RemoveTag_NoteNotFound(t *testing.T) {
	mockNoteRepo := new(MockNoteRepo)
	mockTagRepo := new(MockTagRepo)
	noteTagService := service.NewNoteTagService(mockNoteRepo, mockTagRepo)
	ctx := context.Background()

	userID := uuid.New()
	noteID := uuid.New()
	tagID := uuid.New()

	// Мокируем что заметка не найдена
	mockNoteRepo.On("GetNoteByID", ctx, userID, noteID).Return(nil, nil)

	err := noteTagService.RemoveTag(ctx, userID.String(), noteID.String(), tagID.String())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "note not found")

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteTagService_AddTagsBatch(t *testing.T) {
	mockNoteRepo := new(MockNoteRepo)
	mockTagRepo := new(MockTagRepo)
	noteTagService := service.NewNoteTagService(mockNoteRepo, mockTagRepo)
	ctx := context.Background()

	userID := uuid.New().String()
	noteID1 := uuid.New()
	noteID2 := uuid.New()
	tagID1 := uuid.New()
	tagID2 := uuid.New()

	input := []struct {
		NoteID string
		TagID  string
	}{
		{NoteID: noteID1.String(), TagID: tagID1.String()},
		{NoteID: noteID2.String(), TagID: tagID2.String()},
	}

	expectedNoteTags := []interfaces.NoteTag{
		{NoteID: noteID1, TagID: tagID1},
		{NoteID: noteID2, TagID: tagID2},
	}

	// Мокируем вызов репозитория
	mockNoteRepo.On("AddTagsBatch", ctx, expectedNoteTags).Return(nil)

	err := noteTagService.AddTagsBatch(ctx, userID, input)
	assert.NoError(t, err)

	mockNoteRepo.AssertExpectations(t)
}

func TestNoteTagService_AddTagsBatch_InvalidNoteID(t *testing.T) {
	mockNoteRepo := new(MockNoteRepo)
	mockTagRepo := new(MockTagRepo)
	noteTagService := service.NewNoteTagService(mockNoteRepo, mockTagRepo)
	ctx := context.Background()

	userID := uuid.New().String()

	input := []struct {
		NoteID string
		TagID  string
	}{
		{NoteID: "invalid-uuid", TagID: uuid.New().String()},
	}

	err := noteTagService.AddTagsBatch(ctx, userID, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid noteID in batch")
}

func TestNoteTagService_AddTagsBatch_InvalidTagID(t *testing.T) {
	mockNoteRepo := new(MockNoteRepo)
	mockTagRepo := new(MockTagRepo)
	noteTagService := service.NewNoteTagService(mockNoteRepo, mockTagRepo)
	ctx := context.Background()

	userID := uuid.New().String()

	input := []struct {
		NoteID string
		TagID  string
	}{
		{NoteID: uuid.New().String(), TagID: "invalid-uuid"},
	}

	err := noteTagService.AddTagsBatch(ctx, userID, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid tagID in batch")
}
