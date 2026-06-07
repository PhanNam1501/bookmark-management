package dto

// SuccessResponse is a generic response for successful requests
type SuccessResponse[Data any] struct {
	Message    string       `json:"message,omitempty"`
	Data       Data         `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// ErrorResponse is a response for error cases
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Pagination contains pagination info
type Pagination struct {
	Page  int `json:"page,omitempty" example:"1"`
	Limit int `json:"limit,omitempty" example:"10"`
	Total int `json:"total,omitempty" example:"100"`
}

// NewSuccessResponse creates a success response
func NewSuccessResponse[Data any](message string, data Data) *SuccessResponse[Data] {
	return &SuccessResponse[Data]{
		Message: message,
		Data:    data,
	}
}

// NewSuccessResponseWithPagination creates a success response with pagination
func NewSuccessResponseWithPagination[Data any](message string, data Data, pagination *Pagination) *SuccessResponse[Data] {
	return &SuccessResponse[Data]{
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(code, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// NewErrorResponseWithDetails creates an error response with details
func NewErrorResponseWithDetails(code, message, details string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
}
