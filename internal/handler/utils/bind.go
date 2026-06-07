package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type InputFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type BindError struct {
	Code   string            `json:"code"`
	Errors []InputFieldError `json:"errors"`
}

// BindInputFromRequest binds input from request body with validation
func BindInputFromRequest[T any](c *gin.Context, input *T) error {
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, BindError{
			Code:   "INVALID_INPUT",
			Errors: []InputFieldError{{Field: "body", Message: err.Error()}},
		})
		return err
	}
	return nil
}

// BindInputFromRequestWithReturn binds input from request body and returns the input
func BindInputFromRequestWithReturn[T any](c *gin.Context) (*T, error) {
	input := new(T)
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, BindError{
			Code:   "INVALID_INPUT",
			Errors: []InputFieldError{{Field: "body", Message: err.Error()}},
		})
		return nil, err
	}
	return input, nil
}

// BindInputFromBodyWithoutValidation binds input from body without validation
func BindInputFromBodyWithoutValidation[T any](c *gin.Context, input *T) error {
	if err := c.BindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, BindError{
			Code:   "INVALID_INPUT",
			Errors: []InputFieldError{{Field: "body", Message: err.Error()}},
		})
		return err
	}
	return nil
}

// BindInputFromUri binds input from URI parameters
func BindInputFromUri[T any](c *gin.Context, input *T) error {
	if err := c.ShouldBindUri(input); err != nil {
		c.JSON(http.StatusBadRequest, BindError{
			Code:   "INVALID_INPUT",
			Errors: []InputFieldError{{Field: "uri", Message: err.Error()}},
		})
		return err
	}
	return nil
}

// BindInputFromQuery binds input from query parameters
func BindInputFromQuery[T any](c *gin.Context, input *T) error {
	if err := c.ShouldBindQuery(input); err != nil {
		c.JSON(http.StatusBadRequest, BindError{
			Code:   "INVALID_INPUT",
			Errors: []InputFieldError{{Field: "query", Message: err.Error()}},
		})
		return err
	}
	return nil
}

// BindInput binds from multiple sources (body, uri, query)
func BindInput[T any](c *gin.Context, input *T) error {
	if err := c.ShouldBind(input); err != nil {
		c.JSON(http.StatusBadRequest, BindError{
			Code:   "INVALID_INPUT",
			Errors: []InputFieldError{{Field: "request", Message: err.Error()}},
		})
		return err
	}
	return nil
}

// GetInputWithAuth binds input from request body and validates JWT auth
// Returns input pointer, userID, and error
func GetInputWithAuth[T any](c *gin.Context) (*T, string, error) {
	// Bind input from request body
	input, err := BindInputFromRequestWithReturn[T](c)
	if err != nil {
		return nil, "", err
	}

	// Get user ID from JWT claims
	userID, err := GetUserIDFromRequest(c)
	if err != nil {
		return nil, "", err
	}

	return input, userID, nil
}
