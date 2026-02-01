package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *AppError
		want string
	}{
		{
			name: "without details",
			err:  NewValidationError("test error", nil),
			want: "validation_error: test error",
		},
		{
			name: "with details",
			err:  NewValidationError("test error", map[string]interface{}{"field": "title"}),
			want: "validation_error: test error map[field:title]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("AppError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewValidationError(t *testing.T) {
	details := map[string]interface{}{"field": "title"}
	err := NewValidationError("validation failed", details)

	if err.Code != ErrorCodeValidation {
		t.Errorf("Expected code %s, got %s", ErrorCodeValidation, err.Code)
	}
	if err.Message != "validation failed" {
		t.Errorf("Expected message 'validation failed', got '%s'", err.Message)
	}
	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, err.StatusCode)
	}
	if err.Details == nil {
		t.Error("Expected details to be set")
	}
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("Item")

	if err.Code != ErrorCodeNotFound {
		t.Errorf("Expected code %s, got %s", ErrorCodeNotFound, err.Code)
	}
	if err.Message != "Item not found" {
		t.Errorf("Expected message 'Item not found', got '%s'", err.Message)
	}
	if err.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, err.StatusCode)
	}
}

func TestNewInternalError(t *testing.T) {
	err := NewInternalError("something went wrong")

	if err.Code != ErrorCodeInternal {
		t.Errorf("Expected code %s, got %s", ErrorCodeInternal, err.Code)
	}
	if err.Message != "something went wrong" {
		t.Errorf("Expected message 'something went wrong', got '%s'", err.Message)
	}
	if err.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, err.StatusCode)
	}
}

func TestToErrorResponse(t *testing.T) {
	appErr := NewValidationError("test error", map[string]interface{}{"field": "title"})
	resp := appErr.ToErrorResponse()

	if resp.Error == nil {
		t.Fatal("Expected error detail to be set")
	}
	if resp.Error.Code != ErrorCodeValidation {
		t.Errorf("Expected code %s, got %s", ErrorCodeValidation, resp.Error.Code)
	}
	if resp.Error.Message != "test error" {
		t.Errorf("Expected message 'test error', got '%s'", resp.Error.Message)
	}
	if resp.Error.Details == nil {
		t.Error("Expected details to be set")
	}
}

func TestIsAppError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "is AppError",
			err:  NewValidationError("test", nil),
			want: true,
		},
		{
			name: "is not AppError",
			err:  errors.New("standard error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAppError(tt.err); got != tt.want {
				t.Errorf("IsAppError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want *AppError
	}{
		{
			name: "is AppError",
			err:  NewValidationError("test", nil),
			want: NewValidationError("test", nil),
		},
		{
			name: "is not AppError",
			err:  errors.New("standard error"),
			want: NewInternalError("An unexpected error occurred"),
		},
		{
			name: "nil error",
			err:  nil,
			want: NewInternalError("An unexpected error occurred"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAppError(tt.err)
			if got.Code != tt.want.Code {
				t.Errorf("GetAppError() code = %v, want %v", got.Code, tt.want.Code)
			}
		})
	}
}
