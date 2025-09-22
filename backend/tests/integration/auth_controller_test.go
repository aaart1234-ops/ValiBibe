package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
    "github.com/joho/godotenv"

	"valibibe/internal/repository"
	"valibibe/internal/service"
	"valibibe/internal/router"
	"valibibe/internal/controller"
)

func setupAuthControllerTestRouter(t *testing.T) *gin.Engine {
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
	folderRepo := repository.NewFolderRepo(db)
	folderService := service.NewFolderService(folderRepo)
	folderController := controller.NewFolderController(folderService)
	tagRepo := repository.NewTagRepository(db)
	tagService := service.NewTagService(tagRepo)
	tagController := controller.NewTagController(tagService)

    r := gin.Default()
    router.SetupRoutes(r, tokenService, authController, noteController, folderController, tagController)

    return r
}

func TestRegisterUserHandler_Valid(t *testing.T) {
    r := setupAuthControllerTestRouter(t)

    body := map[string]string{
		"email":    "newuser@example.com",
		"password": "strongpass",
		"nickname": "TestUser",
    }

    jsonValue, _ := json.Marshal(body)

    req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
    assert.Contains(t, w.Body.String(), "newuser@example.com")
}

func TestRegisterUserHandler_Invalid(t *testing.T) {
	r := setupAuthControllerTestRouter(t)

	body := map[string]string{
		"email":    "",
		"password": "",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestLoginUserHandler_Valid(t *testing.T) {
	r := setupAuthControllerTestRouter(t)

	// Сначала зарегистрируем пользователя
	registerBody := map[string]string{
		"email":    "loginuser@example.com",
		"password": "securepass",
		"nickname": "LoginUser",
	}
	registerJSON, _ := json.Marshal(registerBody)
	reqReg, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(registerJSON))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	r.ServeHTTP(wReg, reqReg)
	assert.Equal(t, 200, wReg.Code)

	// Затем логинимся
	loginBody := map[string]string{
		"email":    "loginuser@example.com",
		"password": "securepass",
	}
	loginJSON, _ := json.Marshal(loginBody)
	reqLogin, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginJSON))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	r.ServeHTTP(wLogin, reqLogin)

	assert.Equal(t, 200, wLogin.Code)

	var resp map[string]string
	json.Unmarshal(wLogin.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["token"])
}

func TestLoginUserHandler_Invalid(t *testing.T) {
	r := setupAuthControllerTestRouter(t)

	body := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "wrongpass",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestMeHandler_Authorized(t *testing.T) {
	r := setupAuthControllerTestRouter(t)

	// Зарегистрировать и залогиниться
	registerBody := map[string]string{
		"email":    "me@example.com",
		"password": "mypassword",
		"nickname": "MeUser",
	}
	registerJSON, _ := json.Marshal(registerBody)
	reqReg, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(registerJSON))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	r.ServeHTTP(wReg, reqReg)
	assert.Equal(t, 200, wReg.Code)

	loginBody := map[string]string{
		"email":    "me@example.com",
		"password": "mypassword",
	}
	loginJSON, _ := json.Marshal(loginBody)
	reqLogin, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginJSON))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	r.ServeHTTP(wLogin, reqLogin)

	var resp map[string]string
	json.Unmarshal(wLogin.Body.Bytes(), &resp)
	token := resp["token"]

	// Проверить /auth/me
	reqMe, _ := http.NewRequest("GET", "/auth/me", nil)
	reqMe.Header.Set("Authorization", "Bearer "+token)
	wMe := httptest.NewRecorder()
	r.ServeHTTP(wMe, reqMe)

	assert.Equal(t, 200, wMe.Code)
	assert.Contains(t, wMe.Body.String(), "me@example.com")
}

func TestMeHandler_Unauthorized(t *testing.T) {
	r := setupAuthControllerTestRouter(t)

	req, _ := http.NewRequest("GET", "/auth/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}