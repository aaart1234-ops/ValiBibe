package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"valibibe/internal/controller"
	"valibibe/internal/controller/dto"
	"valibibe/internal/models"
	"valibibe/internal/repository"
	"valibibe/internal/router"
	"valibibe/internal/service"
)

func setupReviewSessionTestRouter(t *testing.T) *gin.Engine {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	db := SetupTestDB(t)

	userRepo := repository.NewUserRepository(db)
	tokenService := service.NewTokenService()
	authService := service.NewAuthService(userRepo, tokenService)
	authController := controller.NewAuthController(authService)

	noteRepo := repository.NewNoteRepository(db)
	noteService := service.NewNoteService(noteRepo)
	assignFolderService := service.NewAssignFolderService(noteRepo)
	noteController := controller.NewNoteController(noteService, assignFolderService)

	folderRepo := repository.NewFolderRepo(db)
	folderService := service.NewFolderService(folderRepo)
	folderController := controller.NewFolderController(folderService)

	tagRepo := repository.NewTagRepository(db)
	tagService := service.NewTagService(tagRepo)
	tagController := controller.NewTagController(tagService)

	noteTagService := service.NewNoteTagService(noteRepo, tagRepo)
	noteTagController := controller.NewNoteTagController(noteTagService)

	reviewSessionService := service.NewReviewSessionService(noteRepo)
	reviewSessionController := controller.NewReviewSessionController(reviewSessionService)

	r := gin.Default()
	router.SetupRoutes(r, tokenService, authController, noteController, folderController, tagController, noteTagController, reviewSessionController)

	return r
}

func TestReviewSession_CreateSession(t *testing.T) {
	r := setupReviewSessionTestRouter(t)

	email := "reviewuser@example.com"
	password := "reviewpass"
	nickname := "ReviewUser"
	token := registerAndLogin(t, r, email, password, nickname)

	// Создаем папку
	folder := createFolder(t, r, token, "Test Folder")

	// Создаем теги
	tag1 := createTag(t, r, token, "Important")
	tag2 := createTag(t, r, token, "Math")

	// Создаем заметки с разными настройками повторения
	nextReview1 := time.Now().Add(-1 * time.Hour)
	nextReview2 := time.Now().Add(-2 * time.Hour)
	note1 := createNoteWithReview(t, r, token, "Note 1", folder.ID.String(), []string{tag1.ID.String()}, 1, &nextReview1)
	note2 := createNoteWithReview(t, r, token, "Note 2", folder.ID.String(), []string{tag2.ID.String()}, 2, &nextReview2)
	_ = createNoteWithReview(t, r, token, "Note 3", "", []string{}, 0, nil) // без повторения

	// Тест 1: Создание сессии без фильтров
	reqBody := dto.ReviewSessionInput{
		Limit: 5,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/review/sessions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response dto.ReviewSessionResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Должны получить заметки готовые к повторению (note1 и note2)
	assert.GreaterOrEqual(t, response.Total, 2)
	assert.LessOrEqual(t, response.Total, 5)

	// Проверяем, что получили заметки готовые к повторению
	noteIDs := make([]string, len(response.Notes))
	for i, note := range response.Notes {
		noteIDs[i] = note.ID
	}

	// Должны получить заметки с next_review_at (note1 и note2)
	// note3 может попасть как случайная заметка при нехватке
	assert.Contains(t, noteIDs, note1.ID.String())
	assert.Contains(t, noteIDs, note2.ID.String())

	// Проверяем, что получили заметки (могут быть как готовые к повторению, так и случайные)
	assert.Greater(t, len(response.Notes), 0, "Should get at least some notes")

	// Тест 2: Создание сессии с фильтром по папке
	folderIDStr := folder.ID.String()
	reqBody = dto.ReviewSessionInput{
		FolderID: &folderIDStr,
		Limit:    10,
	}
	body, _ = json.Marshal(reqBody)
	req, _ = http.NewRequest("POST", "/review/sessions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Все заметки должны быть из указанной папки
	for _, note := range response.Notes {
		assert.Equal(t, folder.ID.String(), note.FolderID)
		assert.Equal(t, folder.Name, note.FolderName)
	}

	// Тест 3: Создание сессии с фильтром по тегам
	reqBody = dto.ReviewSessionInput{
		TagIDs: []string{tag1.ID.String()},
		Limit:  10,
	}
	body, _ = json.Marshal(reqBody)
	req, _ = http.NewRequest("POST", "/review/sessions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Все заметки должны содержать указанный тег
	for _, note := range response.Notes {
		hasTag := false
		for _, tag := range note.Tags {
			if tag.ID == tag1.ID.String() {
				hasTag = true
				break
			}
		}
		assert.True(t, hasTag, "Note should have the specified tag")
	}

	// Тест 4: Валидация лимита
	reqBody = dto.ReviewSessionInput{
		Limit: 150, // больше максимума
	}
	body, _ = json.Marshal(reqBody)
	req, _ = http.NewRequest("POST", "/review/sessions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Лимит должен быть ограничен 100
	assert.LessOrEqual(t, response.Total, 100)
}

// Вспомогательная функция для создания заметки с настройками повторения
func createNoteWithReview(t *testing.T, r *gin.Engine, token, title, folderID string, tagIDs []string, memoryLevel int, nextReviewAt *time.Time) models.Note {
	// Создаем заметку
	body := map[string]interface{}{"title": title, "content": "content"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/notes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	var note models.Note
	json.Unmarshal(w.Body.Bytes(), &note)

	// Назначаем папку если указана
	if folderID != "" {
		assignBody := map[string]string{"folder_id": folderID}
		assignJSON, _ := json.Marshal(assignBody)
		assignReq, _ := http.NewRequest("POST", "/notes/"+note.ID.String()+"/folders", bytes.NewBuffer(assignJSON))
		assignReq.Header.Set("Authorization", "Bearer "+token)
		assignReq.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, assignReq)
		assert.Equal(t, 200, w2.Code)
	}

	// Добавляем теги
	for _, tagID := range tagIDs {
		tagReq, _ := http.NewRequest("POST", "/notes/"+note.ID.String()+"/tags/"+tagID, nil)
		tagReq.Header.Set("Authorization", "Bearer "+token)
		w3 := httptest.NewRecorder()
		r.ServeHTTP(w3, tagReq)
		assert.Equal(t, 200, w3.Code)
	}

	// Обновляем заметку с настройками повторения
	updateBody := map[string]interface{}{
		"title":        title,
		"content":      "content",
		"memory_level": memoryLevel,
	}
	if nextReviewAt != nil {
		updateBody["next_review_at"] = nextReviewAt.Format(time.RFC3339)
	}
	updateJSON, _ := json.Marshal(updateBody)
	updateReq, _ := http.NewRequest("PUT", "/notes/"+note.ID.String(), bytes.NewBuffer(updateJSON))
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateReq.Header.Set("Content-Type", "application/json")
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, updateReq)
	assert.Equal(t, 200, w4.Code)

	// Получаем обновленную заметку
	getReq, _ := http.NewRequest("GET", "/notes/"+note.ID.String(), nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	w5 := httptest.NewRecorder()
	r.ServeHTTP(w5, getReq)
	assert.Equal(t, 200, w5.Code)

	var updatedNote models.Note
	json.Unmarshal(w5.Body.Bytes(), &updatedNote)
	return updatedNote
}
