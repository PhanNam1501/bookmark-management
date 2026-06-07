package bookmark

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type updateBookmarkInput struct {
	URL         string `json:"url" binding:"required,lte=2048"`
	Description string `json:"description" binding:"lte=255"`
}

type UpdateResponse struct {
	Message string `json:"message"`
}

// UpdateBookmark godoc
// @Summary Update bookmark
// @Description Update an existing bookmark (URL and description)
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Bookmark ID"
// @Param request body updateBookmarkInput true "Update data"
// @Success 200 {object} UpdateResponse
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/bookmarks/{id} [put]
func (h *handler) UpdateBookmark(c *gin.Context) {
	claims, _ := c.Get("claims")
	mapClaims := claims.(jwt.MapClaims)
	userID := mapClaims["sub"].(string)

	bookmarkID := c.Param("id")

	var input updateBookmarkInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.s.UpdateBookmark(c, bookmarkID, userID, input.URL, input.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to update bookmark",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UpdateResponse{
		Message: "Success",
	})
}
