package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "my_app_backend/internal/models"
    "my_app_backend/internal/service"
)

type NoteController struct {
    noteService *service.NoteService
}

func NewNoteController(noteService *service.NoteService) *NoteController {
    return &NoteController{noteService: noteService}
}

// CreateNote godoc
// @Summary Создать новую заметку
// @Tags notes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param note body models.NoteInput true "Данные заметки"
// @Success 201 {object} models.Note
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notes [post]
func (c *NoteController) CreateNote(ctx *gin.Context) {
    userID := ctx.MustGet("user_id").(string)

    var input models.NoteInput
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    note, err := c.noteService.CreateNote(ctx, userID, &input)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, note)
}

// GetNoteByID godoc
// @Summary Получить заметку по ID
// @Tags notes
// @Security BearerAuth
// @Produce json
// @Param id path string true "Note ID"
// @Success 200 {object} models.Note
// @Failure 404 {object} map[string]string
// @Router /notes/{id} [get]
func (c *NoteController) GetNoteByID(ctx *gin.Context) {
    userID := ctx.MustGet("user_id").(string)
    id := ctx.Param("id")

    note, err := c.noteService.GetNoteByID(ctx, userID, id)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, note)
}

// GetAllNotes godoc
// @Summary Получить заметки пользователя с фильтрацией и сортировкой
// @Description Возвращает список заметок текущего пользователя. Поддерживаются параметры поиска и сортировки.
// @Tags notes
// @Security BearerAuth
// @Produce json
// @Param search query string false "Поиск по заголовку или содержимому"
// @Param sort_by query string false "Поле сортировки: created_at (дата создания), next_review_at (дата следующего повторения)" Enums(created_at, next_review_at) default(created_at)
// @Param order query string false "Порядок сортировки: asc (по возрастанию), desc (по убыванию)" Enums(asc, desc) default(desc)
// @Success 200 {array} models.Note
// @Failure 500 {object} map[string]string
// @Router /notes [get]
func (c *NoteController) GetAllNotes(ctx *gin.Context) {
    userID := ctx.MustGet("user_id").(string)

    filter := models.NoteFilter{
        UserID: userID,
        Search: ctx.Query("search"),
        SortBy: ctx.DefaultQuery("sort_by", "created_at"),
        Order:  ctx.DefaultQuery("order", "desc"),
    }

    notes, err := c.noteService.GetAllNotesByUserID(ctx, &filter)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, notes)
}

// UpdateNote godoc
// @Summary Обновить заметку
// @Tags notes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Note ID"
// @Param note body models.NoteInput true "Обновлённые данные заметки"
// @Success 200 {object} models.Note
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notes/{id} [put]
func (c *NoteController) UpdateNote(ctx *gin.Context) {
    userID := ctx.MustGet("user_id").(string)
    id := ctx.Param("id")

    var input models.NoteInput
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    note, err := c.noteService.UpdateNote(ctx, userID, id, &input)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, note)
}

// ArchiveNote godoc
// @Summary Архивировать заметку
// @Tags notes
// @Security BearerAuth
// @Produce json
// @Param id path string true "Note ID"
// @Success 200 {object} models.Note
// @Failure 500 {object} map[string]string
// @Router /notes/{id}/archive [post]
func (c *NoteController) ArchiveNote(ctx *gin.Context) {
    userID := ctx.MustGet("user_id").(string)
    id := ctx.Param("id")

    updatedNote, err := c.noteService.ArchiveNote(ctx, userID, id)  // Теперь возвращает (*models.Note, error)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, updatedNote)  // Возвращаем саму заметку
}

// DeleteNote godoc
// @Summary Удалить заметку
// @Tags notes
// @Security BearerAuth
// @Produce json
// @Param id path string true "Note ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notes/{id} [delete]
func (c *NoteController) DeleteNote(ctx *gin.Context) {
    userID := ctx.MustGet("user_id").(string)
    id := ctx.Param("id")

    err := c.noteService.DeleteNote(ctx, userID, id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
}

// ReviewNoteHandler godoc
// @Summary Обновить память (review) по заметке
// @Tags notes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Note ID"
// @Param input body models.ReviewInput true "Вспомнил или нет"
// @Success 200 {object} models.Note
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notes/{id}/review [post]
func (c *NoteController) ReviewNoteHandler(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	noteID := ctx.Param("id")

	var input models.ReviewInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := c.noteService.UpdateMemoryLevel(ctx, userID, noteID, input.Remembered)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    // Получаем обновлённую заметку и возвращаем её
    updatedNote, err := c.noteService.GetNoteByID(ctx, userID, noteID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated note"})
        return
    }

    ctx.JSON(http.StatusOK, updatedNote) // Теперь в ответе всегда будет memoryLevel (0 или увеличенный)
}
