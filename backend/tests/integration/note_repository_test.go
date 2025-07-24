package integration

import (
	"context"
	"testing"
	"fmt"

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

	// --- GetNoteByID ---
	gotNote, err := noteRepo.GetNoteByID(ctx, note.ID.String())
	require.NoError(t, err)
	assert.Equal(t, note.Title, gotNote.Title)

	// --- GetNoteByIDAndUserID (positive) ---
	gotNoteByUser, err := noteRepo.GetNoteByIDAndUserID(ctx, note.ID.String(), user.ID.String())
	require.NoError(t, err)
	assert.Equal(t, note.ID, gotNoteByUser.ID)

	// --- GetNoteByIDAndUserID (wrong user) ---
	gotNoteWrongUser, err := noteRepo.GetNoteByIDAndUserID(ctx, note.ID.String(), uuid.New().String())
	require.NoError(t, err)
	assert.Nil(t, gotNoteWrongUser)

	// --- Bulk create notes ---
	for i := 1; i <= 5; i++ {
		n := &models.Note{
			ID:      uuid.New(),
			UserID:  user.ID,
			Title:   fmt.Sprintf("Note %d", i),
			Content: "bulk content",
		}
		require.NoError(t, noteRepo.CreateNote(ctx, n))
	}

	// --- GetAllNotesByUserID: filter, sort, pagination ---
	filter := &models.NoteFilter{
		UserID: user.ID.String(),
		SortBy: "created_at",
		Order:  "asc",
		Limit:  3,
		Offset: 0,
	}
	paginated, err := noteRepo.GetAllNotesByUserID(ctx, filter)
	require.NoError(t, err)
	assert.Equal(t, int64(6), paginated.Total)
	assert.Len(t, paginated.Notes, 3)

	// --- Search ---
	searchFilter := &models.NoteFilter{
		UserID: user.ID.String(),
		Search: "note 3", // case-insensitive
	}
	searchResult, err := noteRepo.GetAllNotesByUserID(ctx, searchFilter)
	require.NoError(t, err)
	assert.Len(t, searchResult.Notes, 1)
	assert.Equal(t, "Note 3", searchResult.Notes[0].Title)

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

	// --- UnarchiveNote ---
    err = noteRepo.UnArchiveNote(ctx, note.ID.String())
    require.NoError(t, err)

    unarchivedNote, err := noteRepo.GetNoteByID(ctx, note.ID.String())
    require.NoError(t, err)
    assert.False(t, unarchivedNote.Archived)

    filterAfterUnarchive := &models.NoteFilter{UserID: user.ID.String()}
    notesAfterUnarchive, err := noteRepo.GetAllNotesByUserID(ctx, filterAfterUnarchive)
    require.NoError(t, err)

    	found := false
    	for _, n := range notesAfterUnarchive.Notes {
    		if n.ID == note.ID {
    			found = true
    			break
    		}
    	}
    	assert.True(t, found, "Unarchived note should be returned in GetAllNotesByUserID")


	// --- GetAllNotesByUserID: проверка игнорирования archived ---
	filterAfterArchive := &models.NoteFilter{
    	UserID:   user.ID.String(),
    	Archived: ptr(false), // Фильтрация только по неархивным
    }
	notesAfterArchive, err := noteRepo.GetAllNotesByUserID(ctx, filterAfterArchive)
	require.NoError(t, err)
	for _, n := range notesAfterArchive.Notes {
		assert.False(t, n.Archived)
	}

	// --- DeleteNote ---
	err = noteRepo.DeleteNote(ctx, note.ID.String())
	require.NoError(t, err)

	deletedNote, err := noteRepo.GetNoteByID(ctx, note.ID.String())
	require.NoError(t, err)
	assert.Nil(t, deletedNote)
}

func ptr[T any](v T) *T {
	return &v
}

