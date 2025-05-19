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
    userID := ctx.MustGet("userID").(string)

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
    userID := ctx.MustGet("userID").(string)
    id := ctx.Param("id")

    note, err := c.noteService.GetNoteByID(ctx, userID, id)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, note)
}

func (c *NoteController) GetAllNotes(ctx *gin.Context) {
    userID := ctx.MustGet("userID").(string)

    notes, err := c.noteService.GetAllNotesByUserID(ctx, userID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, notes)
}

func (c *NoteController) UpdateNote(ctx *gin.Context) {
    userID := ctx.MustGet("userID").(string)
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
    userID := ctx.MustGet("userID").(string)
    id := ctx.Param("id")

    err := c.noteService.ArchiveNote(ctx, userID, id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Note archived"})
}

func (c *NoteController) DeleteNote(ctx *gin.Context) {
    userID := ctx.MustGet("userID").(string)
    id := ctx.Param("id")

    err := c.noteService.DeleteNote(ctx, userID, id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
}
