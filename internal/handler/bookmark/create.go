package bookmark

import (
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/handler/dto"
	"github.com/PhanNam1501/bookmark-management/internal/handler/utils"
	"github.com/PhanNam1501/bookmark-management/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type addBookmarkInput struct {
	Description string `json:"description" example:"Google" validate:"lte=255"`
	URL         string `json:"url" example:"https://www.google.com" validate:"required,url,lte=2048"`
}

type CreateBookmarkResponse struct {
	Message string         `json:"message"`
	Data    *model.Bookmark `json:"data"`
}

// CreateBookmark godoc
// @Summary Create a new bookmark
// @Description Create a new bookmark for the authenticated user
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body addBookmarkInput true "Bookmark request"
// @Success 200 {object} CreateBookmarkResponse
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/bookmarks [post]
func (h *handler) CreateBookmark(c *gin.Context) {
	input, uid, err := utils.GetInputWithAuth[addBookmarkInput](c)
	if err != nil {
		return
	}

	res, err := h.s.CreateBookmark(c, uid, input.URL, input.Description)
	if err != nil {
		log.Error().Err(err).Str("userId", uid).Msg("CreateBookmark error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create bookmark",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &dto.SuccessResponse[*model.Bookmark]{
		Message: "Bookmark created successfully",
		Data:    res,
	})
}
