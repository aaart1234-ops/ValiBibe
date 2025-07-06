package integration

import (
    "fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/joho/godotenv"

	"my_app_backend/internal/models"
	"my_app_backend/internal/repository"
	"my_app_backend/internal/service"
	"my_app_backend/internal/router"
	"my_app_backend/internal/controller"
)

func setupNoteControllerTestRouter(t *testing.T) *gin.Engine {
	err := godotenv.Load("../../../.env")
	assert.NoError(t, err)

	db := setupTestDB(t)

	userRepo := repository.NewUserRepository(db)
	tokenService := service.NewTokenService()
	authService := service.NewAuthService(userRepo, tokenService)
	authController := controller.NewAuthController(authService)

	noteRepo := repository.NewNoteRepository(db)
	noteService := service.NewNoteService(noteRepo)
	noteController := controller.NewNoteController(noteService)

	r := gin.Default()
	router.SetupRoutes(r, tokenService, authController, noteController)

	return r
}

// helper: регистрация пользователя и получение токена
func registerAndLogin(t *testing.T, r *gin.Engine, email, password, nickname string) string {
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

func TestNoteCRUDFlow(t *testing.T) {
	r := setupNoteControllerTestRouter(t)

	email := "noteuser@example.com"
	password := "noteuserpass"
	nickname := "NoteUser"

	token := registerAndLogin(t, r, email, password, nickname)

	// --- Create Note ---
	createBody := map[string]string{
		"title":   "Test Note Title",
		"content": "This is the test note content",
	}
	createJSON, _ := json.Marshal(createBody)
	reqCreate, _ := http.NewRequest("POST", "/notes", bytes.NewBuffer(createJSON))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	wCreate := httptest.NewRecorder()
	r.ServeHTTP(wCreate, reqCreate)
	assert.Equal(t, 201, wCreate.Code)

	var createdNote models.Note
	err := json.Unmarshal(wCreate.Body.Bytes(), &createdNote)
	require.NoError(t, err)
	assert.Equal(t, createBody["title"], createdNote.Title)
	assert.Equal(t, createBody["content"], createdNote.Content)
	assert.NotEmpty(t, createdNote.ID)

	// --- Get Note by ID ---
	reqGet, _ := http.NewRequest("GET", "/notes/"+createdNote.ID.String(), nil)
	reqGet.Header.Set("Authorization", "Bearer "+token)
	wGet := httptest.NewRecorder()
	r.ServeHTTP(wGet, reqGet)
	assert.Equal(t, 200, wGet.Code)

	var gotNote models.Note
	err = json.Unmarshal(wGet.Body.Bytes(), &gotNote)
	require.NoError(t, err)
	assert.Equal(t, createdNote.ID, gotNote.ID)
	assert.Equal(t, createdNote.Title, gotNote.Title)

	// --- Update Note ---
	updateBody := map[string]string{
		"title":   "Updated Note Title",
		"content": "Updated note content",
	}
	updateJSON, _ := json.Marshal(updateBody)
	reqUpdate, _ := http.NewRequest("PUT", "/notes/"+createdNote.ID.String(), bytes.NewBuffer(updateJSON))
	reqUpdate.Header.Set("Content-Type", "application/json")
	reqUpdate.Header.Set("Authorization", "Bearer "+token)
	wUpdate := httptest.NewRecorder()
	r.ServeHTTP(wUpdate, reqUpdate)
	assert.Equal(t, 200, wUpdate.Code)

	var updatedNote models.Note
	err = json.Unmarshal(wUpdate.Body.Bytes(), &updatedNote)
	require.NoError(t, err)
	assert.Equal(t, updateBody["title"], updatedNote.Title)
	assert.Equal(t, updateBody["content"], updatedNote.Content)

	// --- Archive Note ---
	reqArchive, _ := http.NewRequest("POST", "/notes/"+createdNote.ID.String()+"/archive", nil)
	reqArchive.Header.Set("Authorization", "Bearer "+token)
	wArchive := httptest.NewRecorder()
	r.ServeHTTP(wArchive, reqArchive)
	assert.Equal(t, 200, wArchive.Code)

	var archivedNote models.Note
	err = json.Unmarshal(wArchive.Body.Bytes(), &archivedNote)
	require.NoError(t, err)
	assert.True(t, archivedNote.Archived)

	// --- Delete Note ---
	reqDelete, _ := http.NewRequest("DELETE", "/notes/"+createdNote.ID.String(), nil)
	reqDelete.Header.Set("Authorization", "Bearer "+token)
	wDelete := httptest.NewRecorder()
	r.ServeHTTP(wDelete, reqDelete)
	assert.Equal(t, 200, wDelete.Code)

	// Проверяем, что после удаления запрос по id возвращает 404 или пустое тело
	reqGetDeleted, _ := http.NewRequest("GET", "/notes/"+createdNote.ID.String(), nil)
	reqGetDeleted.Header.Set("Authorization", "Bearer "+token)
	wGetDeleted := httptest.NewRecorder()
	r.ServeHTTP(wGetDeleted, reqGetDeleted)
	assert.Equal(t, 404, wGetDeleted.Code)
}

func TestNotesListAndUnauthorized(t *testing.T) {
	r := setupNoteControllerTestRouter(t)

	email := "listuser@example.com"
	password := "listuserpass"
	nickname := "ListUser"

	token := registerAndLogin(t, r, email, password, nickname)

	// Создаем несколько заметок
	for i := 1; i <= 3; i++ {
		createBody := map[string]string{
			"title":   "Note " + string(rune(i+'0')),
			"content": "Content for note " + string(rune(i+'0')),
		}
		createJSON, _ := json.Marshal(createBody)
		reqCreate, _ := http.NewRequest("POST", "/notes", bytes.NewBuffer(createJSON))
		reqCreate.Header.Set("Content-Type", "application/json")
		reqCreate.Header.Set("Authorization", "Bearer "+token)
		wCreate := httptest.NewRecorder()
		r.ServeHTTP(wCreate, reqCreate)
		assert.Equal(t, 201, wCreate.Code)
	}

	// --- Get All Notes ---
	reqGetAll, _ := http.NewRequest("GET", "/notes", nil)
	reqGetAll.Header.Set("Authorization", "Bearer "+token)
	wGetAll := httptest.NewRecorder()
	r.ServeHTTP(wGetAll, reqGetAll)
	assert.Equal(t, 200, wGetAll.Code)

	var paginated struct {
        Notes []models.Note `json:"notes"`
    }
    err := json.Unmarshal(wGetAll.Body.Bytes(), &paginated)
    require.NoError(t, err)
    assert.GreaterOrEqual(t, len(paginated.Notes), 3)

	// --- Unauthorized access ---
	reqUnauthorized, _ := http.NewRequest("GET", "/notes", nil)
	wUnauthorized := httptest.NewRecorder()
	r.ServeHTTP(wUnauthorized, reqUnauthorized)
	assert.Equal(t, 401, wUnauthorized.Code)
}

func TestReviewNoteHandler(t *testing.T) {
	r := setupAuthControllerTestRouter(t)

	// Регистрация и логин
	registerBody := map[string]string{
		"email":    "review@example.com",
		"password": "pass1234",
		"nickname": "ReviewUser",
	}
	registerJSON, _ := json.Marshal(registerBody)
	reqReg, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(registerJSON))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	r.ServeHTTP(wReg, reqReg)
	assert.Equal(t, 200, wReg.Code)

	loginBody := map[string]string{
		"email":    "review@example.com",
		"password": "pass1234",
	}
	loginJSON, _ := json.Marshal(loginBody)
	reqLogin, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginJSON))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	r.ServeHTTP(wLogin, reqLogin)
	assert.Equal(t, 200, wLogin.Code)

	var loginResp map[string]string
	json.Unmarshal(wLogin.Body.Bytes(), &loginResp)
	token := loginResp["token"]

	// Создание заметки
	noteBody := map[string]string{
		"title":   "Memory Note",
		"content": "Content to be remembered",
	}
	noteJSON, _ := json.Marshal(noteBody)
	reqNote, _ := http.NewRequest("POST", "/notes", bytes.NewBuffer(noteJSON))
	reqNote.Header.Set("Authorization", "Bearer "+token)
	reqNote.Header.Set("Content-Type", "application/json")
	wNote := httptest.NewRecorder()
	r.ServeHTTP(wNote, reqNote)
	assert.Equal(t, 201, wNote.Code)

	var noteResp map[string]interface{}
	json.Unmarshal(wNote.Body.Bytes(), &noteResp)
	noteID := noteResp["id"].(string)

	// --- Case 1: Remembered = true ---
	review := models.ReviewInput{Remembered: true}
	reviewJSON, _ := json.Marshal(review)
	reqReview, _ := http.NewRequest("POST", "/notes/"+noteID+"/review", bytes.NewBuffer(reviewJSON))
	reqReview.Header.Set("Authorization", "Bearer "+token)
	reqReview.Header.Set("Content-Type", "application/json")
	wReview := httptest.NewRecorder()
	r.ServeHTTP(wReview, reqReview)
	fmt.Println("Review response body:", wReview.Body.String()) // Логирование ответа
	assert.Equal(t, 200, wReview.Code)

	// Проверяем, что memoryLevel > 0 и nextReviewAt задан
	reqGet, _ := http.NewRequest("GET", "/notes/"+noteID, nil)
	reqGet.Header.Set("Authorization", "Bearer "+token)
	wGet := httptest.NewRecorder()
	r.ServeHTTP(wGet, reqGet)
	fmt.Println("GET response body:", wGet.Body.String())

	var note map[string]interface{}
	json.Unmarshal(wGet.Body.Bytes(), &note)
	fmt.Printf("Parsed note: %+v\n", note) // Логирование структуры

	memoryLevel, ok := note["memoryLevel"].(float64)
    if !ok {
        assert.Fail(t, "memoryLevel is missing or has wrong type")
    } else {
        assert.Greater(t, int(memoryLevel), 0)
    }

    nextReviewAt, exists := note["next_review_at"]
    if !exists {
        assert.Fail(t, "nextReviewAt field is missing")
    } else if nextReviewAt == nil {
        assert.Fail(t, "nextReviewAt should not be null after positive review")
    } else {
        assert.NotEmpty(t, nextReviewAt)
    }

	// --- Case 2: Remembered = false ---
	reviewFalse := models.ReviewInput{Remembered: false}
	reviewFalseJSON, _ := json.Marshal(reviewFalse)
	reqReviewFalse, _ := http.NewRequest("POST", "/notes/"+noteID+"/review", bytes.NewBuffer(reviewFalseJSON))
	reqReviewFalse.Header.Set("Authorization", "Bearer "+token)
	reqReviewFalse.Header.Set("Content-Type", "application/json")
	wReviewFalse := httptest.NewRecorder()
	r.ServeHTTP(wReviewFalse, reqReviewFalse)
	assert.Equal(t, 200, wReviewFalse.Code)

	// Проверяем, что memoryLevel сброшен и nextReviewAt == null
	reqGet2, _ := http.NewRequest("GET", "/notes/"+noteID, nil)
	reqGet2.Header.Set("Authorization", "Bearer "+token)
	wGet2 := httptest.NewRecorder()
	r.ServeHTTP(wGet2, reqGet2)

	var noteAfter map[string]interface{}
	json.Unmarshal(wGet2.Body.Bytes(), &noteAfter)
	assert.Equal(t, float64(0), noteAfter["memoryLevel"])
	assert.Nil(t, noteAfter["nextReviewAt"])
}
