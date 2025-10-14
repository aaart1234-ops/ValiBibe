package controller

import (
	"net/http"

	"valibibe/internal/controller/dto"
	"valibibe/internal/service"

	"github.com/gin-gonic/gin"
)

type FolderController struct {
	folderService *service.FolderService
}

func NewFolderController(folderService *service.FolderService) *FolderController {
	return &FolderController{folderService: folderService}
}

// CreateFolder godoc
// @Summary Создать папку
// @Tags folders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param folder body dto.FolderCreateInput true "Данные папки"
// @Success 201 {object} models.Folder
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /folders [post]
func (c *FolderController) CreateFolder(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	var input dto.FolderCreateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder, err := c.folderService.CreateFolder(ctx, userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, folder)
}

// GetFolderTree godoc
// @Summary Получить дерево папок пользователя
// @Tags folders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.FolderNode
// @Failure 500 {object} map[string]string
// @Router /folders/tree [get]
func (c *FolderController) GetFolderTree(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	tree, err := c.folderService.GetFolderTree(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tree)
}

// UpdateFolder godoc
// @Summary Обновить папку
// @Tags folders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Folder ID"
// @Param folder body dto.FolderUpdateInput true "Обновлённые данные папки"
// @Success 200 {object} models.Folder
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /folders/{id} [put]
func (c *FolderController) UpdateFolder(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	id := ctx.Param("id")

	var input dto.FolderUpdateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder, err := c.folderService.UpdateFolder(ctx, userID, id, input)
	if err != nil {
		if err.Error() == "folder not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, folder)
}

// DeleteFolder godoc
// @Summary Удалить папку (с каскадом заметок и подпапок)
// @Tags folders
// @Security BearerAuth
// @Param id path string true "Folder ID"
// @Success 204 {string} string "No Content"
// @Failure 500 {object} map[string]string
// @Router /folders/{id} [delete]
func (c *FolderController) DeleteFolder(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	id := ctx.Param("id")

	if err := c.folderService.DeleteFolder(ctx, userID, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
