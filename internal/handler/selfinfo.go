package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// GetCurrentUser godoc
// @Summary      Get current user information
// @Description  Get the authenticated user's information from JWT token
// @Tags         users
// @Security     Bearer
// @Produce      json
// @Success      200 {object} model.User
// @Failure      401 {object} map[string]string "Unauthorized or invalid token"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /v1/self/info [get]
func (h *userHandler) GetCurrentUser(c *gin.Context) {
	values, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, ok := values.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	uid, ok := claims["sub"].(string)
	if !ok || uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	user, err := h.svc.GetCurrentUser(c, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error"})
		return
	}

	c.JSON(http.StatusOK, user)
}
