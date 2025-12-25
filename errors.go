package mailbreeze

import (
	"fmt"
	"net/http"
)

// Error represents an API error.
type Error struct {
	// StatusCode is the HTTP status code.
	StatusCode int

	// Code is the machine-readable error code.
	Code string

	// Message is the human-readable error message.
	Message string

	// RequestID is the unique request ID for debugging.
	RequestID string

	// RetryAfter is the number of seconds to wait before retrying (for rate limits).
	RetryAfter int

	// Details contains additional error details.
	Details map[string]interface{}
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("mailbreeze: %s (code: %s, status: %d, request_id: %s)", e.Message, e.Code, e.StatusCode, e.RequestID)
	}
	return fmt.Sprintf("mailbreeze: %s (code: %s, status: %d)", e.Message, e.Code, e.StatusCode)
}

// newError creates a new Error.
func newError(statusCode int, message, code, requestID string, retryAfter int, details map[string]interface{}) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		RequestID:  requestID,
		RetryAfter: retryAfter,
		Details:    details,
	}
}

// newErrorFromStatus creates an Error from an HTTP status code.
func newErrorFromStatus(statusCode int, message, code, requestID string, retryAfter int) *Error {
	if code == "" {
		code = codeFromStatus(statusCode)
	}
	return newError(statusCode, message, code, requestID, retryAfter, nil)
}

// codeFromStatus returns a default error code for an HTTP status.
func codeFromStatus(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "VALIDATION_ERROR"
	case http.StatusUnauthorized:
		return "AUTHENTICATION_ERROR"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusTooManyRequests:
		return "RATE_LIMIT_EXCEEDED"
	default:
		if statusCode >= 500 {
			return "SERVER_ERROR"
		}
		return "UNKNOWN_ERROR"
	}
}

// IsAuthenticationError returns true if the error is an authentication error.
func IsAuthenticationError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == http.StatusUnauthorized
	}
	return false
}

// IsValidationError returns true if the error is a validation error.
func IsValidationError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == http.StatusBadRequest
	}
	return false
}

// IsNotFoundError returns true if the error is a not found error.
func IsNotFoundError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == http.StatusNotFound
	}
	return false
}

// IsRateLimitError returns true if the error is a rate limit error.
func IsRateLimitError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// IsServerError returns true if the error is a server error.
func IsServerError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode >= 500
	}
	return false
}

// GetRetryAfter returns the retry-after duration in seconds, or 0 if not applicable.
func GetRetryAfter(err error) int {
	if e, ok := err.(*Error); ok {
		return e.RetryAfter
	}
	return 0
}
