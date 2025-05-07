package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"my_app_backend/internal/models"
	"my_app_backend/internal/repository"
)

func TestNoteRepository(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	noteRepo := repository.NewNoteRepository(db)

	// Создание пользователя
	user := &models.User{
		ID:           uuid.New(),
		Nickname:     "TestUser",
		Email:        "test@example.com",
		PasswordHash: "hashed-password",
	}
	require.NoError(t, db.Create(user).Error)

	// --- CreateNote ---
	note := &models.Note{
		ID:      uuid.New(),
		UserID:  user.ID,
		Title:   "Test Note",
		Content: "This is a test note",
	}

	err := noteRepo.CreateNote(ctx, note)
	require.NoError(t, err)
	assert.NotEmpty(t, note.CreatedAt)

	// --- GetNoteByID ---
	gotNote, err := noteRepo.GetNoteByID(ctx, note.ID.String())
	require.NoError(t, err)
	assert.Equal(t, note.Title, gotNote.Title)

	// --- GetAllNotesByUserID ---
	notes, err := noteRepo.GetAllNotesByUserID(ctx, user.ID.String())
	require.NoError(t, err)
	assert.Len(t, notes, 1)
	assert.Equal(t, "Test Note", notes[0].Title)

	// --- UpdateNote ---
	note.Title = "Updated Title"
	err = noteRepo.UpdateNote(ctx, note)
	require.NoError(t, err)

	updatedNote, err := noteRepo.GetNoteByID(ctx, note.ID.String())
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedNote.Title)

	// --- ArchiveNote ---
	err = noteRepo.ArchiveNote(ctx, note.ID.String())
	require.NoError(t, err)

	archivedNote, err := noteRepo.GetNoteByID(ctx, note.ID.String())
	require.NoError(t, err)
	assert.True(t, archivedNote.Archived)

	// --- DeleteNote ---
	err = noteRepo.DeleteNote(ctx, note.ID.String())
	require.NoError(t, err)

	deletedNote, err := noteRepo.GetNoteByID(ctx, note.ID.String())
	require.NoError(t, err)
	assert.Nil(t, deletedNote)
}
