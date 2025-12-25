package mailbreeze

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPClientErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   map[string]interface{}
		checkErr   func(error) bool
	}{
		{
			name:       "400 validation error",
			statusCode: http.StatusBadRequest,
			response: map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    "INVALID_EMAIL",
					"message": "Invalid email format",
				},
			},
			checkErr: IsValidationError,
		},
		{
			name:       "401 authentication error",
			statusCode: http.StatusUnauthorized,
			response: map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    "AUTHENTICATION_ERROR",
					"message": "Invalid API key",
				},
			},
			checkErr: IsAuthenticationError,
		},
		{
			name:       "404 not found error",
			statusCode: http.StatusNotFound,
			response: map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    "NOT_FOUND",
					"message": "Email not found",
				},
			},
			checkErr: IsNotFoundError,
		},
		{
			name:       "429 rate limit error",
			statusCode: http.StatusTooManyRequests,
			response: map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests",
				},
			},
			checkErr: IsRateLimitError,
		},
		{
			name:       "500 server error",
			statusCode: http.StatusInternalServerError,
			response: map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    "SERVER_ERROR",
					"message": "Internal server error",
				},
			},
			checkErr: IsServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client := NewClient("sk_test_123", WithBaseURL(server.URL), WithMaxRetries(0))

			_, err := client.Emails.Get(context.Background(), "test")

			if err == nil {
				t.Fatal("expected error")
			}

			if !tt.checkErr(err) {
				t.Errorf("expected %s check to return true, got error: %v", tt.name, err)
			}
		})
	}
}

func TestHTTPClientRetryAfterHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "Too many requests",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL), WithMaxRetries(0))

	_, err := client.Emails.Get(context.Background(), "test")

	if err == nil {
		t.Fatal("expected error")
	}

	retryAfter := GetRetryAfter(err)
	if retryAfter != 60 {
		t.Errorf("expected retry-after 60, got %d", retryAfter)
	}
}

func TestHTTPClientRequestID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "req_abc123")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL), WithMaxRetries(0))

	_, err := client.Emails.Get(context.Background(), "test")

	if err == nil {
		t.Fatal("expected error")
	}

	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatal("expected *Error type")
	}

	if apiErr.RequestID != "req_abc123" {
		t.Errorf("expected request ID 'req_abc123', got '%s'", apiErr.RequestID)
	}
}

func TestHTTPClient204NoContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	err := client.Lists.Delete(context.Background(), "list_123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHTTPClientRetry(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   map[string]interface{}{"message": "Error"},
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":         "email_123",
				"from":       "a@example.com",
				"to":         []string{"b@example.com"},
				"status":     "delivered",
				"created_at": "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123",
		WithBaseURL(server.URL),
		WithMaxRetries(3),
		WithHTTPClient(&http.Client{Timeout: 100 * time.Millisecond}),
	)

	email, err := client.Emails.Get(context.Background(), "email_123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}

	if email.ID != "email_123" {
		t.Errorf("expected email ID 'email_123', got '%s'", email.ID)
	}
}

func TestHTTPClientNoRetryOn400(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   map[string]interface{}{"message": "Invalid"},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL), WithMaxRetries(3))

	_, err := client.Emails.Get(context.Background(), "test")

	if err == nil {
		t.Fatal("expected error")
	}

	if attempts != 1 {
		t.Errorf("expected 1 attempt (no retry), got %d", attempts)
	}
}

func TestHTTPClientUserAgent(t *testing.T) {
	var userAgent string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"sent": 0, "delivered": 0, "bounced": 0, "complained": 0,
				"opened": 0, "clicked": 0, "unsubscribed": 0,
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	client.Emails.Stats(context.Background())

	if userAgent != "mailbreeze-go/"+Version {
		t.Errorf("expected User-Agent 'mailbreeze-go/%s', got '%s'", Version, userAgent)
	}
}

func TestAPIKeyRedactedInDebugOutput(t *testing.T) {
	secretKey := "sk_live_super_secret_api_key_12345"
	client := NewClient(secretKey)

	// Test Client.String()
	clientStr := client.String()
	if contains(clientStr, secretKey) {
		t.Error("API key should not appear in Client.String() output")
	}
	if !contains(clientStr, "[REDACTED]") {
		t.Error("Client.String() should show [REDACTED]")
	}

	// Test HTTPClient.String()
	httpClientStr := client.httpClient.String()
	if contains(httpClientStr, secretKey) {
		t.Error("API key should not appear in HTTPClient.String() output")
	}
	if !contains(httpClientStr, "[REDACTED]") {
		t.Error("HTTPClient.String() should show [REDACTED]")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
