package types

import "net/http"

// APIError represents an API error
type APIError struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

// Error returns the error message
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new API error
func NewAPIError(statusCode, code int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// Common errors
var (
	ErrBadRequest       = NewAPIError(http.StatusBadRequest, 400, "Bad request")
	ErrUnauthorized     = NewAPIError(http.StatusUnauthorized, 401, "Unauthorized")
	ErrForbidden        = NewAPIError(http.StatusForbidden, 403, "Forbidden")
	ErrNotFound         = NewAPIError(http.StatusNotFound, 404, "Not found")
	ErrTooManyRequests  = NewAPIError(http.StatusTooManyRequests, 429, "Too many requests")
	ErrInternalServerError = NewAPIError(http.StatusInternalServerError, 500, "Internal server error")
)
