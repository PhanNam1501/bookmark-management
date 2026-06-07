package bookmark

import (
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/handler/dto"
	"github.com/PhanNam1501/bookmark-management/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type getBookmarksInput struct {
	Page  int `form:"page" validate:"gte=1"`
	Limit int `form:"limit" validate:"gte=1"`
}

type GetBookmarksResponse struct {
	Message    string                `json:"message"`
	Data       []*model.Bookmark     `json:"data"`
	Pagination *dto.Pagination       `json:"pagination,omitempty"`
}

// GetBookmarks godoc
// @Summary Get user bookmarks
// @Description Get paginated list of bookmarks for the authenticated user
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "Page number (starting from 1)" default(1)
// @Param limit query int true "Items per page" default(10)
// @Success 200 {object} GetBookmarksResponse "List of bookmarks"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/bookmarks [get]
func (h *handler) GetBookmarks(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	uid, exists := mapClaims["sub"]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var input getBookmarksInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid query parameters",
		})
		return
	}

	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}

	res, err := h.s.GetBookmarks(c, uid.(string), input.Limit, input.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get bookmarks",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &GetBookmarksResponse{
		Message: "Bookmarks retrieved successfully",
		Data:    res.Bookmarks,
		Pagination: &dto.Pagination{
			Page:  input.Page,
			Limit: input.Limit,
			Total: int(res.Count),
		},
	})
}
