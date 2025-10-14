package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"valibibe/internal/controller"
	"valibibe/internal/repository"
	"valibibe/internal/router"
	"valibibe/internal/service"
)

func setupNoteTagControllerTestRouter(t *testing.T) *gin.Engine {
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

	r := gin.Default()
	router.SetupRoutes(r, tokenService, authController, noteController, folderController, tagController, noteTagController, nil)

	return r
}

// helper: регистрация пользователя и получение токена
func registerAndLoginForNoteTag(t *testing.T, r *gin.Engine, email, password, nickname string) string {
	// регистрация
	registerBody := map[string]string{
		"email":    email,
		"password": password,
		"nickname": nickname,
	}
	registerJSON, _ := json.Marshal(registerBody)
	reqReg, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(registerJSON))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	r.ServeHTTP(wReg, reqReg)
	require.Equal(t, 200, wReg.Code)

	// логин
	loginBody := map[string]string{
		"email":    email,
		"password": password,
	}
	loginJSON, _ := json.Marshal(loginBody)
	reqLogin, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginJSON))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	r.ServeHTTP(wLogin, reqLogin)
	require.Equal(t, 200, wLogin.Code)

	var resp map[string]string
	err := json.Unmarshal(wLogin.Body.Bytes(), &resp)
	require.NoError(t, err)
	token, ok := resp["token"]
	require.True(t, ok)
	return token
}

func TestNoteTagController_AddTag(t *testing.T) {
	r := setupNoteTagControllerTestRouter(t)

	email := "notetag@example.com"
	password := "notetagpass"
	nickname := "NoteTagUser"

	token := registerAndLoginForNoteTag(t, r, email, password, nickname)

	// Создаем заметку
	createBody := map[string]string{
		"title":   "Test Note for Tags",
		"content": "This is a test note for tag operations",
	}
	createJSON, _ := json.Marshal(createBody)
	reqCreate, _ := http.NewRequest("POST", "/notes", bytes.NewBuffer(createJSON))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	wCreate := httptest.NewRecorder()
	r.ServeHTTP(wCreate, reqCreate)
	assert.Equal(t, 201, wCreate.Code)

	var createdNote map[string]interface{}
	err := json.Unmarshal(wCreate.Body.Bytes(), &createdNote)
	require.NoError(t, err)
	noteID := createdNote["id"].(string)

	// Создаем тег
	tagBody := map[string]string{
		"name": "Test Tag",
	}
	tagJSON, _ := json.Marshal(tagBody)
	reqTag, _ := http.NewRequest("POST", "/tags", bytes.NewBuffer(tagJSON))
	reqTag.Header.Set("Content-Type", "application/json")
	reqTag.Header.Set("Authorization", "Bearer "+token)
	wTag := httptest.NewRecorder()
	r.ServeHTTP(wTag, reqTag)
	assert.Equal(t, 201, wTag.Code)

	var createdTag map[string]interface{}
	err = json.Unmarshal(wTag.Body.Bytes(), &createdTag)
	require.NoError(t, err)
	tagID := createdTag["id"].(string)

	// Добавляем тег к заметке
	reqAddTag, _ := http.NewRequest("POST", "/notes/"+noteID+"/tags/"+tagID, nil)
	reqAddTag.Header.Set("Authorization", "Bearer "+token)
	wAddTag := httptest.NewRecorder()
	r.ServeHTTP(wAddTag, reqAddTag)
	assert.Equal(t, 200, wAddTag.Code)
}

func TestNoteTagController_RemoveTag(t *testing.T) {
	r := setupNoteTagControllerTestRouter(t)

	email := "notetag2@example.com"
	password := "notetag2pass"
	nickname := "NoteTagUser2"

	token := registerAndLoginForNoteTag(t, r, email, password, nickname)

	// Создаем заметку
	createBody := map[string]string{
		"title":   "Test Note for Tag Removal",
		"content": "This is a test note for tag removal",
	}
	createJSON, _ := json.Marshal(createBody)
	reqCreate, _ := http.NewRequest("POST", "/notes", bytes.NewBuffer(createJSON))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	wCreate := httptest.NewRecorder()
	r.ServeHTTP(wCreate, reqCreate)
	assert.Equal(t, 201, wCreate.Code)

	var createdNote map[string]interface{}
	err := json.Unmarshal(wCreate.Body.Bytes(), &createdNote)
	require.NoError(t, err)
	noteID := createdNote["id"].(string)

	// Создаем тег
	tagBody := map[string]string{
		"name": "Test Tag for Removal",
	}
	tagJSON, _ := json.Marshal(tagBody)
	reqTag, _ := http.NewRequest("POST", "/tags", bytes.NewBuffer(tagJSON))
	reqTag.Header.Set("Content-Type", "application/json")
	reqTag.Header.Set("Authorization", "Bearer "+token)
	wTag := httptest.NewRecorder()
	r.ServeHTTP(wTag, reqTag)
	assert.Equal(t, 201, wTag.Code)

	var createdTag map[string]interface{}
	err = json.Unmarshal(wTag.Body.Bytes(), &createdTag)
	require.NoError(t, err)
	tagID := createdTag["id"].(string)

	// Добавляем тег к заметке
	reqAddTag, _ := http.NewRequest("POST", "/notes/"+noteID+"/tags/"+tagID, nil)
	reqAddTag.Header.Set("Authorization", "Bearer "+token)
	wAddTag := httptest.NewRecorder()
	r.ServeHTTP(wAddTag, reqAddTag)
	assert.Equal(t, 200, wAddTag.Code)

	// Удаляем тег у заметки
	reqRemoveTag, _ := http.NewRequest("DELETE", "/notes/"+noteID+"/tags/"+tagID, nil)
	reqRemoveTag.Header.Set("Authorization", "Bearer "+token)
	wRemoveTag := httptest.NewRecorder()
	r.ServeHTTP(wRemoveTag, reqRemoveTag)
	assert.Equal(t, 200, wRemoveTag.Code)
}

func TestNoteTagController_AddTagsBatch(t *testing.T) {
	r := setupNoteTagControllerTestRouter(t)

	email := "notetag3@example.com"
	password := "notetag3pass"
	nickname := "NoteTagUser3"

	token := registerAndLoginForNoteTag(t, r, email, password, nickname)

	// Создаем несколько заметок
	var noteIDs []string
	for i := 1; i <= 3; i++ {
		createBody := map[string]string{
			"title":   "Test Note " + string(rune(i+'0')),
			"content": "This is test note " + string(rune(i+'0')),
		}
		createJSON, _ := json.Marshal(createBody)
		reqCreate, _ := http.NewRequest("POST", "/notes", bytes.NewBuffer(createJSON))
		reqCreate.Header.Set("Content-Type", "application/json")
		reqCreate.Header.Set("Authorization", "Bearer "+token)
		wCreate := httptest.NewRecorder()
		r.ServeHTTP(wCreate, reqCreate)
		assert.Equal(t, 201, wCreate.Code)

		var createdNote map[string]interface{}
		err := json.Unmarshal(wCreate.Body.Bytes(), &createdNote)
		require.NoError(t, err)
		noteIDs = append(noteIDs, createdNote["id"].(string))
	}

	// Создаем несколько тегов
	var tagIDs []string
	for i := 1; i <= 2; i++ {
		tagBody := map[string]string{
			"name": "Test Tag " + string(rune(i+'0')),
		}
		tagJSON, _ := json.Marshal(tagBody)
		reqTag, _ := http.NewRequest("POST", "/tags", bytes.NewBuffer(tagJSON))
		reqTag.Header.Set("Content-Type", "application/json")
		reqTag.Header.Set("Authorization", "Bearer "+token)
		wTag := httptest.NewRecorder()
		r.ServeHTTP(wTag, reqTag)
		assert.Equal(t, 201, wTag.Code)

		var createdTag map[string]interface{}
		err := json.Unmarshal(wTag.Body.Bytes(), &createdTag)
		require.NoError(t, err)
		tagIDs = append(tagIDs, createdTag["id"].(string))
	}

	// Массовое добавление тегов к заметкам
	batchBody := []map[string]string{
		{"note_id": noteIDs[0], "tag_id": tagIDs[0]},
		{"note_id": noteIDs[1], "tag_id": tagIDs[0]},
		{"note_id": noteIDs[2], "tag_id": tagIDs[1]},
	}
	batchJSON, _ := json.Marshal(batchBody)
	reqBatch, _ := http.NewRequest("POST", "/notes/tags/batch", bytes.NewBuffer(batchJSON))
	reqBatch.Header.Set("Content-Type", "application/json")
	reqBatch.Header.Set("Authorization", "Bearer "+token)
	wBatch := httptest.NewRecorder()
	r.ServeHTTP(wBatch, reqBatch)
	assert.Equal(t, 200, wBatch.Code)
}

func TestNoteTagController_Unauthorized(t *testing.T) {
	r := setupNoteTagControllerTestRouter(t)

	// Попытка добавить тег без авторизации
	req, _ := http.NewRequest("POST", "/notes/some-note-id/tags/some-tag-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// Попытка массового добавления тегов без авторизации
	batchBody := []map[string]string{
		{"note_id": "some-note-id", "tag_id": "some-tag-id"},
	}
	batchJSON, _ := json.Marshal(batchBody)
	reqBatch, _ := http.NewRequest("POST", "/notes/tags/batch", bytes.NewBuffer(batchJSON))
	reqBatch.Header.Set("Content-Type", "application/json")
	wBatch := httptest.NewRecorder()
	r.ServeHTTP(wBatch, reqBatch)
	assert.Equal(t, 401, wBatch.Code)
}
