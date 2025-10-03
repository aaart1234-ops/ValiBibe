package controller

import (
	"net/http"
	"strconv"

	"valibibe/internal/controller/dto"
	"valibibe/internal/service"

	"github.com/gin-gonic/gin"
)

type NoteController struct {
	noteService         *service.NoteService
	assignFolderService *service.AssignFolderService
}

func NewNoteController(
	noteService *service.NoteService,
	assignFolderService *service.AssignFolderService,
) *NoteController {
	return &NoteController{
		noteService:         noteService,
		assignFolderService: assignFolderService,
	}
}

// CreateNote godoc
// @Summary Создать новую заметку
// @Tags notes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param note body dto.NoteInput true "Данные заметки"
// @Success 201 {object} models.Note
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notes [post]
func (c *NoteController) CreateNote(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	var input dto.NoteInput
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
// @Param limit query int false "Максимальное количество записей" minimum(1) default(10)
// @Param offset query int false "Смещение для пагинации" minimum(0) default(0)
// @Param folder_id query string false "ID папки для фильтрации заметок по папке"
// @Param tag_ids query []string false "Массив ID тегов для фильтрации заметок по тегам (через tag_ids[]=id1&tag_ids[]=id2)"
// @Success 200 {object} dto.PaginatedNotes
// @Failure 500 {object} map[string]string
// @Router /notes [get]
func (c *NoteController) GetAllNotes(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	archivedStr := ctx.Query("archived")
	var archived *bool
	if archivedStr != "" {
		parsed, err := strconv.ParseBool(archivedStr)
		if err == nil {
			archived = &parsed
		}
	}
	folderID := ctx.Query("folder_id")
	var folderIDPtr *string
	if folderID != "" {
		folderIDPtr = &folderID
	}
	tagIDs := ctx.QueryArray("tag_ids[]")

	filter := dto.NoteFilter{
		UserID:   userID,
		Search:   ctx.Query("search"),
		SortBy:   ctx.DefaultQuery("sort_by", "created_at"),
		Order:    ctx.DefaultQuery("order", "desc"),
		Limit:    limit,
		Offset:   offset,
		Archived: archived,
		FolderID: folderIDPtr,
		TagIDs:   tagIDs,
	}

	result, err := c.noteService.GetAllNotesByUserID(ctx, &filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// UpdateNote godoc
// @Summary Обновить заметку
// @Tags notes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Note ID"
// @Param note body dto.NoteInput true "Обновлённые данные заметки"
// @Success 200 {object} models.Note
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notes/{id} [put]
func (c *NoteController) UpdateNote(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	id := ctx.Param("id")

	var input dto.NoteInput
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

	updatedNote, err := c.noteService.ArchiveNote(ctx, userID, id) // Теперь возвращает (*models.Note, error)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedNote) // Возвращаем саму заметку
}

// UnArchiveNote godoc
// @Summary Разархивировать заметку
// @Tags notes
// @Security BearerAuth
// @Produce json
// @Param id path string true "Note ID"
// @Success 200 {object} models.Note
// @Failure 500 {object} map[string]string
// @Router /notes/{id}/unarchive [post]
func (c *NoteController) UnArchiveNote(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	id := ctx.Param("id")

	updatedNote, err := c.noteService.UnArchiveNote(ctx, userID, id) // Теперь возвращает (*models.Note, error)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedNote) // Возвращаем саму заметку
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
// @Param input body dto.ReviewInput true "Вспомнил или нет"
// @Success 200 {object} models.Note
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notes/{id}/review [post]
func (c *NoteController) ReviewNoteHandler(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	noteID := ctx.Param("id")

	var input dto.ReviewInput
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

// AssignFolder godoc
// @Summary      Assign a folder to a note
// @Description  Link a note to a folder (each note can belong to only one folder).
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id path string true "Note ID" example("a12b34cd-5678-90ef-1234-567890abcdef")
// @Param        input body dto.AssignFolderInput true "Folder assignment input"
// @Success      200 {string} string "ok"
// @Failure      400 {object} map[string]string
// @Router       /notes/{id}/folders [post]
func (c *NoteController) AssignFolder(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	noteID := ctx.Param("id")

	var input dto.AssignFolderInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.assignFolderService.AssignFolder(ctx, userID, noteID, input.FolderID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// RemoveFolder godoc
// @Summary      Remove folder from a note
// @Description  Unlink a note from its folder.
// @Tags         notes
// @Produce      json
// @Param        id path string true "Note ID" example("a12b34cd-5678-90ef-1234-567890abcdef")
// @Param        folderId path string true "Folder ID" example("d290f1ee-6c54-4b01-90e6-d701748f0851")
// @Success      200 {string} string "ok"
// @Failure      400 {object} map[string]string
// @Router       /notes/{id}/folders/{folderId} [delete]
func (c *NoteController) RemoveFolder(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	noteID := ctx.Param("id")
	folderID := ctx.Param("folderId")

	if err := c.assignFolderService.RemoveFolder(ctx, userID, noteID, folderID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// BatchAssignFolder godoc
// @Summary      Batch assign or remove folder for multiple notes
// @Description  Mass update notes to link or unlink them from a folder.
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        input body dto.BatchAssignFolderInput true "Batch folder assignment input"
// @Success      200 {string} string "ok"
// @Failure      400 {object} map[string]string
// @Router       /notes/batch/folders [post]
func (c *NoteController) BatchAssignFolder(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	var input dto.BatchAssignFolderInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.assignFolderService.BatchAssignFolder(ctx, userID, input.NoteIDs, input.FolderID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
