package handler

import (
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type RedirectURL interface {
	RedirectURL(c *gin.Context)
}

type redirectURLHandler struct {
	svc service.URLRedirect
}

func NewRedirectURLHandler(svc service.URLRedirect) RedirectURL {
	return &redirectURLHandler{
		svc: svc,
	}
}

type redirectURLRequest struct {
	Code string `uri:"code" binding:"required"`
}

type RedirectURLResponse struct {
	URL string `json:"url"`
}

// RedirectURL redirects the short code to its original long URL.
// @Summary      Get Original URL
// @Description  Retrieve the original long URL associated with the provided short code from the path.
// @Tags         Redirect
// @Accept       json
// @Produce      json
// @Param        code   path      string  true  "The unique short code (e.g., abc123X)"
// @Success      200    {object}  RedirectURLResponse "Successfully retrieved the URL"
// @Failure      400    {object}  map[string]string   "Invalid or missing code in path"
// @Failure      500    {object}  map[string]string   "Internal server error"
// @Router       /v1/links/redirect/{code} [get]
func (r *redirectURLHandler) RedirectURL(c *gin.Context) {
	input := &redirectURLRequest{}
	if err := c.ShouldBindUri(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := r.svc.GetRedirectURL(c, input.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server err"})
		return
	}

	c.JSON(http.StatusOK, RedirectURLResponse{
		URL: url,
	})
}
