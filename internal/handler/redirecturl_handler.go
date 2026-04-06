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

// RedirectURL redirects the short code to its original long URL.
// @Summary      Redirect to Original URL
// @Description  Retrieves the original long URL and performs a 301 redirect.
// @Tags         Redirect
// @Param        code   path      string  true  "The unique short code"
// @Success      301    {string}  string  "Redirecting to the original URL"
// @Failure      400    {object}  map[string]string  "Invalid or missing code"
// @Failure      404    {object}  map[string]string  "Short code not found"
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

	c.Redirect(http.StatusMovedPermanently, url)
}
