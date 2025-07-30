package errors

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AppError represents an application error
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	StatusCode int    `json:"status_code"`
	Timestamp  string `json:"timestamp"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// New creates a new AppError
func New(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// Common error codes
const (
	ErrCodeInvalidRequest      = "INVALID_REQUEST"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeForbidden           = "FORBIDDEN"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeInternalError       = "INTERNAL_ERROR"
	ErrCodeValidationFailed    = "VALIDATION_FAILED"
	ErrCodeResourceNotFound    = "RESOURCE_NOT_FOUND"
	ErrCodePermissionDenied    = "PERMISSION_DENIED"
	ErrCodeBadRequest          = "BAD_REQUEST"
	ErrCodeConflict            = "CONFLICT"
	ErrCodeUnprocessableEntity = "UNPROCESSABLE_ENTITY"
)

// Common errors
var (
	ErrInvalidRequest      = New(ErrCodeInvalidRequest, "Invalid request", http.StatusBadRequest)
	ErrUnauthorized        = New(ErrCodeUnauthorized, "Unauthorized", http.StatusUnauthorized)
	ErrForbidden           = New(ErrCodeForbidden, "Forbidden", http.StatusForbidden)
	ErrNotFound            = New(ErrCodeNotFound, "Resource not found", http.StatusNotFound)
	ErrInternalError       = New(ErrCodeInternalError, "Internal server error", http.StatusInternalServerError)
	ErrValidationFailed    = New(ErrCodeValidationFailed, "Validation failed", http.StatusBadRequest)
	ErrResourceNotFound    = New(ErrCodeResourceNotFound, "Resource not found", http.StatusNotFound)
	ErrPermissionDenied    = New(ErrCodePermissionDenied, "Permission denied", http.StatusForbidden)
	ErrBadRequest          = New(ErrCodeBadRequest, "Bad request", http.StatusBadRequest)
	ErrConflict            = New(ErrCodeConflict, "Conflict", http.StatusConflict)
	ErrUnprocessableEntity = New(ErrCodeUnprocessableEntity, "Unprocessable entity", http.StatusUnprocessableEntity)
)

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, err error) {
	var appErr *AppError

	// Check if it's already an AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	} else {
		// Convert to AppError
		appErr = ErrInternalError.WithDetails(err.Error())
	}

	c.JSON(appErr.StatusCode, appErr)
}

// RespondWithValidationError sends a validation error response
func RespondWithValidationError(c *gin.Context, details string) {
	err := ErrValidationFailed.WithDetails(details)
	c.JSON(err.StatusCode, err)
}

// RespondWithPermissionError sends a permission denied error response
func RespondWithPermissionError(c *gin.Context, resource string) {
	err := ErrPermissionDenied.WithDetails("Insufficient permissions for " + resource)
	c.JSON(err.StatusCode, err)
}

// RespondWithNotFoundError sends a not found error response
func RespondWithNotFoundError(c *gin.Context, resource string) {
	err := ErrResourceNotFound.WithDetails(resource + " not found")
	c.JSON(err.StatusCode, err)
}
