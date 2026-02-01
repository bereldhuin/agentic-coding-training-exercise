package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents an error code
type ErrorCode string

const (
	ErrorCodeValidation ErrorCode = "validation_error"
	ErrorCodeNotFound   ErrorCode = "not_found"
	ErrorCodeInternal   ErrorCode = "internal_error"
	ErrorCodeBadRequest ErrorCode = "bad_request"
	ErrorCodeConflict   ErrorCode = "conflict"
)

// AppError represents an application error
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	StatusCode int                    `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != nil && len(e.Details) > 0 {
		return fmt.Sprintf("%s: %s %v", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error (for errors.Is/As)
func (e *AppError) Unwrap() error {
	return nil
}

// NewValidationError creates a new validation error
func NewValidationError(message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:       ErrorCodeValidation,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusBadRequest,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       ErrorCodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

// NewInternalError creates a new internal error
func NewInternalError(message string) *AppError {
	return &AppError{
		Code:       ErrorCodeInternal,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:       ErrorCodeBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:       ErrorCodeConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// WrapError wraps an error with context
func WrapError(err error, code ErrorCode, message string, statusCode int) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// ErrorResponse represents the canonical error response format
type ErrorResponse struct {
	Error *ErrorDetail `json:"error"`
}

// ErrorDetail represents error details
type ErrorDetail struct {
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ToErrorResponse converts an AppError to ErrorResponse
func (e *AppError) ToErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		Error: &ErrorDetail{
			Code:    e.Code,
			Message: e.Message,
			Details: e.Details,
		},
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError extracts an AppError from an error, or creates a new internal error
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return NewInternalError("An unexpected error occurred")
}
