package main

import (
	"fmt"
	"net/http"
)

// Custom error types for better error handling

// APIError represents errors from API requests
type APIError struct {
	StatusCode int
	Message    string
	URL        string
}

func (e APIError) Error() string {
	return fmt.Sprintf("API request failed: %s (status: %d, url: %s)", e.Message, e.StatusCode, e.URL)
}

// IsNotFound checks if the error is a 404 Not Found
func (e APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// NewAPIError creates a new APIError
func NewAPIError(statusCode int, message, url string) APIError {
	return APIError{
		StatusCode: statusCode,
		Message:    message,
		URL:        url,
	}
}

// ValidationError represents input validation errors
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s: %s", e.Field, e.Message)
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// JSONError represents JSON parsing errors
type JSONError struct {
	Message string
	Data    string
}

func (e JSONError) Error() string {
	return fmt.Sprintf("JSON parsing error: %s", e.Message)
}

// NewJSONError creates a new JSONError
func NewJSONError(message, data string) JSONError {
	return JSONError{
		Message: message,
		Data:    data,
	}
}
