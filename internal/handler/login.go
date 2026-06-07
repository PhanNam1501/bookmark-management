package handler

import (
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/handler/dto"
	"github.com/gin-gonic/gin"
)

type loginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=8"`
}

type LoginTokenResponse struct {
	Token string `json:"token"`
}

// LoginResponse is the response for login endpoint
type LoginResponse struct {
	Message string             `json:"message"`
	Data    LoginTokenResponse `json:"data"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user with username and password, return JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param request body loginInput true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 401 {object} dto.ErrorResponse "Invalid username or password"
// @Router /users/login [post]
func (h *userHandler) Login(c *gin.Context) {
	input := &loginInput{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_REQUEST", err.Error()))
		return
	}

	token, err := h.svc.Login(c, input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("AUTH_FAILED", "Invalid username or password"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Login successful", LoginTokenResponse{
		Token: token,
	}))
}

type bodyInput struct {
	Name string `json:"name" binding:"required"`
}

type queryInput struct {
	Page  int `uri:"page" binding:"required"`
	Limit int `uri:"limit" binding:"required"`
}

type headerInput struct {
	Count string `header:"count" binding:"required"`
}

func (h *userHandler) Test(c *gin.Context) {
	code := c.Param("code")
	println(code)

	input := &bodyInput{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	println(input.Name)

	queryInput := &queryInput{}
	if err := c.ShouldBindQuery(queryInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query input"})
		return
	}

	println(queryInput.Page, queryInput.Limit)

	hInput := &headerInput{}
	if err := c.ShouldBindHeader(hInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid header"})
		return
	}
	println(hInput.Count)
}
