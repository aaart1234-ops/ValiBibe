package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"valibibe/internal/controller/dto"
	"valibibe/internal/service"
)

type ReviewSessionController struct {
	reviewSessionService *service.ReviewSessionService
}

func NewReviewSessionController(reviewSessionService *service.ReviewSessionService) *ReviewSessionController {
	return &ReviewSessionController{
		reviewSessionService: reviewSessionService,
	}
}

// CreateReviewSession godoc
// @Summary Создать сессию повторения
// @Description Создает сессию повторения с фильтрацией по папке и тегам. Возвращает заметки готовые к повторению, при нехватке добавляет случайные заметки.
// @Tags review-sessions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body dto.ReviewSessionInput true "Параметры сессии повторения"
// @Success 200 {object} dto.ReviewSessionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /review/sessions [post]
func (c *ReviewSessionController) CreateReviewSession(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	var input dto.ReviewSessionInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	result, err := c.reviewSessionService.CreateReviewSession(ctx, userID, &input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

