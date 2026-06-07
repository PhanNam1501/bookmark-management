package bookmark

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type DeleteResponse struct {
	Message string `json:"message"`
}

// DeleteBookmark godoc
// @Summary Delete bookmark
// @Description Delete a bookmark by ID (must be owned by current user)
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Bookmark ID"
// @Success 200 {object} DeleteResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/bookmarks/{id} [delete]
func (h *handler) DeleteBookmark(c *gin.Context) {
	claims, _ := c.Get("claims")
	mapClaims := claims.(jwt.MapClaims)
	userID := mapClaims["sub"].(string)

	bookmarkID := c.Param("id")

	err := h.s.DeleteBookmark(c, bookmarkID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to delete bookmark",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, DeleteResponse{
		Message: "Success",
	})
}
