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

func (c *NoteController) GetAllNotes(ctx *gin.Context) {
    userID := ctx.MustGet("user_id").(string)

    notes, err := c.noteService.GetAllNotesByUserID(ctx, userID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, notes)
}

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
