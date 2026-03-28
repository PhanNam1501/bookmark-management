package handler

import (
	"net/http"

	"github.com/PhanNam1501/bookmark-management/internal/service"
	"github.com/gin-gonic/gin"
)

type ShortenURL interface {
	ShortenURL(c *gin.Context)
}

type shortenURLHandler struct {
	svc service.ShortenURL
}

func NewShortenURLHandler(svc service.ShortenURL) ShortenURL {
	return &shortenURLHandler{
		svc: svc,
	}
}

type shortenURLRequest struct {
	URL string `json:"url"`
}

type ShortenURLResponse struct {
	Key string `json:"key"`
}

// @Summary Shorten URL
// @Tags ShortenURL
// @Accept json
// @Produce json
// @Param request body shortenURLRequest true "Request body"
// @Success 200 {object} map[string]string
// @Router /shorten [post]
func (s *shortenURLHandler) ShortenURL(c *gin.Context) {
	input := &shortenURLRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	key, err := s.svc.ShortenURL(c, input.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server err"})
		return
	}

	c.JSON(http.StatusOK, ShortenURLResponse{
		Key: key,
	})
}
