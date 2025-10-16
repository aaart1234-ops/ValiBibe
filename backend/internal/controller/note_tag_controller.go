package controller

import (
	"errors"
	"net/http"
	"valibibe/internal/controller/dto"
	apperrors "valibibe/internal/errors"
	"valibibe/internal/service"

	"github.com/gin-gonic/gin"
)

type NoteTagController struct {
	noteTagService *service.NoteTagService
}

func NewNoteTagController(noteTagService *service.NoteTagService) *NoteTagController {
	return &NoteTagController{
		noteTagService: noteTagService,
	}
}

// AddTag godoc
// @Summary Добавить тег к заметке
// @Tags note-tags
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Note ID" example("a12b34cd-5678-90ef-1234-567890abcdef")
// @Param tagId path string true "Tag ID" example("b23c45de-6789-01fg-2345-678901bcdefg")
// @Success 200 {string} string "ok"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /notes/{id}/tags/{tagId} [post]
func (c *NoteTagController) AddTag(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	noteID := ctx.Param("id")
	tagID := ctx.Param("tagId")

	if err := c.noteTagService.AddTag(ctx, userID, noteID, tagID); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Note or tag not found"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// RemoveTag godoc
// @Summary Удалить тег у заметки
// @Tags note-tags
// @Security BearerAuth
// @Produce json
// @Param id path string true "Note ID" example("a12b34cd-5678-90ef-1234-567890abcdef")
// @Param tagId path string true "Tag ID" example("b23c45de-6789-01fg-2345-678901bcdefg")
// @Success 200 {string} string "ok"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /notes/{id}/tags/{tagId} [delete]
func (c *NoteTagController) RemoveTag(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	noteID := ctx.Param("id")
	tagID := ctx.Param("tagId")

	if err := c.noteTagService.RemoveTag(ctx, userID, noteID, tagID); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Note or tag not found"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// AddTagsBatch godoc
// @Summary Массовое добавление тегов к заметкам
// @Tags note-tags
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body []dto.NoteTagInput true "Массив связей заметка-тег"
// @Success 200 {string} string "ok"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /notes/tags/batch [post]
func (c *NoteTagController) AddTagsBatch(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	var input []dto.NoteTagInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.noteTagService.AddTagsBatch(ctx, userID, input); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "One or more notes/tags not found"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
