package handler

import (
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type User interface {
	RegisterUser(c *gin.Context)
	Login(c *gin.Context)
	GetCurrentUser(c *gin.Context)
	Test(g *gin.Context)
}

type userHandler struct {
	svc service.User
}

func NewUser(svc service.User) User {
	return &userHandler{svc: svc}
}

type redisterInputBody struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Create a new user account with username, password, display name, and email
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input body redisterInputBody true "User registration data"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string "Invalid input"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users/register [post]
func (u *userHandler) RegisterUser(c *gin.Context) {
	input := &redisterInputBody{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	res, err := u.svc.CreateUser(c, input.Username, input.Password, input.DisplayName, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, res)
}
