package handler

import (
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type LinkShortURL interface {
	LinkShortenURL(c *gin.Context)
}

type linkshortURLHandler struct {
	svc service.ShortenURL
}

func NewLinkShortenHandler(svc service.ShortenURL) LinkShortURL {
	return &linkshortURLHandler{
		svc: svc,
	}
}

type linkShortenURLRequest struct {
	ExpTime int    `json:"exp"`
	URL     string `json:"url"`
}

type LinkShortenURLResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// @Summary Link Shorten URL
// @Tags LinkShortenURL
// @Accept json
// @Produce json
// @Param request body linkShortenURLRequest true "Request body"
// @Success 200 {object} map[string]string
// @Router /v1/links/shorten [post]
func (s *linkshortURLHandler) LinkShortenURL(c *gin.Context) {
	input := &linkShortenURLRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	key, err := s.svc.LinkShortenURL(c, input.URL, input.ExpTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server err"})
		return
	}

	c.JSON(http.StatusOK, LinkShortenURLResponse{
		Code:    key,
		Message: "Shorten URL generated successfully!",
	})
}
