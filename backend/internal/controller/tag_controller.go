package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"valibibe/internal/controller/dto"
	apperrors "valibibe/internal/errors"
	"valibibe/internal/service"
)

type TagController struct {
	tagService *service.TagService
}

func NewTagController(tagService *service.TagService) *TagController {
	return &TagController{tagService: tagService}
}

// CreateTag godoc
// @Summary      Создать тег
// @Description  Создаёт новый тег для пользователя
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        body  body      dto.TagCreateInput  true  "Данные тега"
// @Success      201   {object}  models.Tag
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /tags [post]
func (tc *TagController) CreateTag(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	var input dto.TagCreateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	tag, err := tc.tagService.CreateTag(ctx, userID, input)
	if err != nil {
		if err.Error() == "tag with this name already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, tag)
}

// ListTags godoc
// @Summary      Список тегов
// @Description  Получает список всех тегов пользователя
// @Tags         tags
// @Produce      json
// @Success      200  {array}   models.Tag
// @Failure      500  {object}  map[string]string
// @Router       /tags [get]
func (tc *TagController) ListTags(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	tags, err := tc.tagService.GetTags(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tags)
}

// UpdateTag godoc
// @Summary      Обновить тег
// @Description  Обновляет имя тега по ID
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id    path      string             true  "ID тега"
// @Param        body  body      dto.TagUpdateInput true  "Новые данные"
// @Success      200   {object}  models.Tag
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /tags/{id} [put]
func (tc *TagController) UpdateTag(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	tagID := ctx.Param("id")

	var input dto.TagUpdateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	tag, err := tc.tagService.UpdateTag(ctx, userID, tagID, input)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
			return
		}
		if err.Error() == "tag with this name already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, tag)
}

// DeleteTag godoc
// @Summary      Удалить тег
// @Description  Удаляет тег по ID
// @Tags         tags
// @Produce      json
// @Param        id   path      string  true  "ID тега"
// @Success      204  "Тег удалён"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tags/{id} [delete]
func (tc *TagController) DeleteTag(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	tagID := ctx.Param("id")

	if err := tc.tagService.DeleteTag(ctx, userID, tagID); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
