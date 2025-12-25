package mailbreeze

import (
	"net/http"
	"testing"
)

func TestErrorMessage(t *testing.T) {
	err := &Error{
		StatusCode: 400,
		Code:       "VALIDATION_ERROR",
		Message:    "Invalid email format",
		RequestID:  "req_123",
	}

	expected := "mailbreeze: Invalid email format (code: VALIDATION_ERROR, status: 400, request_id: req_123)"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestErrorMessageWithoutRequestID(t *testing.T) {
	err := &Error{
		StatusCode: 400,
		Code:       "VALIDATION_ERROR",
		Message:    "Invalid email format",
	}

	expected := "mailbreeze: Invalid email format (code: VALIDATION_ERROR, status: 400)"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestIsAuthenticationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "authentication error",
			err:      &Error{StatusCode: http.StatusUnauthorized},
			expected: true,
		},
		{
			name:     "other error",
			err:      &Error{StatusCode: http.StatusBadRequest},
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAuthenticationError(tt.err); got != tt.expected {
				t.Errorf("IsAuthenticationError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	if !IsValidationError(&Error{StatusCode: http.StatusBadRequest}) {
		t.Error("expected true for 400 status")
	}
	if IsValidationError(&Error{StatusCode: http.StatusNotFound}) {
		t.Error("expected false for 404 status")
	}
}

func TestIsNotFoundError(t *testing.T) {
	if !IsNotFoundError(&Error{StatusCode: http.StatusNotFound}) {
		t.Error("expected true for 404 status")
	}
	if IsNotFoundError(&Error{StatusCode: http.StatusBadRequest}) {
		t.Error("expected false for 400 status")
	}
}

func TestIsRateLimitError(t *testing.T) {
	if !IsRateLimitError(&Error{StatusCode: http.StatusTooManyRequests}) {
		t.Error("expected true for 429 status")
	}
	if IsRateLimitError(&Error{StatusCode: http.StatusBadRequest}) {
		t.Error("expected false for 400 status")
	}
}

func TestIsServerError(t *testing.T) {
	if !IsServerError(&Error{StatusCode: http.StatusInternalServerError}) {
		t.Error("expected true for 500 status")
	}
	if !IsServerError(&Error{StatusCode: 503}) {
		t.Error("expected true for 503 status")
	}
	if IsServerError(&Error{StatusCode: http.StatusBadRequest}) {
		t.Error("expected false for 400 status")
	}
}

func TestGetRetryAfter(t *testing.T) {
	err := &Error{RetryAfter: 60}
	if got := GetRetryAfter(err); got != 60 {
		t.Errorf("expected 60, got %d", got)
	}

	if got := GetRetryAfter(nil); got != 0 {
		t.Errorf("expected 0 for nil, got %d", got)
	}
}

func TestCodeFromStatus(t *testing.T) {
	tests := []struct {
		status   int
		expected string
	}{
		{http.StatusBadRequest, "VALIDATION_ERROR"},
		{http.StatusUnauthorized, "AUTHENTICATION_ERROR"},
		{http.StatusForbidden, "FORBIDDEN"},
		{http.StatusNotFound, "NOT_FOUND"},
		{http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED"},
		{http.StatusInternalServerError, "SERVER_ERROR"},
		{http.StatusServiceUnavailable, "SERVER_ERROR"},
		{418, "UNKNOWN_ERROR"}, // I'm a teapot
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := codeFromStatus(tt.status); got != tt.expected {
				t.Errorf("codeFromStatus(%d) = %s, want %s", tt.status, got, tt.expected)
			}
		})
	}
}
