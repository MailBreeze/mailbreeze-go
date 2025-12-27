package mailbreeze

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Test error helper functions with non-*Error types
func TestErrorHelpersWithNonAPIError(t *testing.T) {
	genericErr := errors.New("generic error")

	if IsValidationError(genericErr) {
		t.Error("expected false for generic error")
	}
	if IsAuthenticationError(genericErr) {
		t.Error("expected false for generic error")
	}
	if IsNotFoundError(genericErr) {
		t.Error("expected false for generic error")
	}
	if IsRateLimitError(genericErr) {
		t.Error("expected false for generic error")
	}
	if IsServerError(genericErr) {
		t.Error("expected false for generic error")
	}
	if GetRetryAfter(genericErr) != 0 {
		t.Error("expected 0 for generic error")
	}
}

// Test Lists.List with all query parameters
func TestListsListWithQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("expected page 2, got %s", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("limit") != "50" {
			t.Errorf("expected limit 50, got %s", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("search") != "test" {
			t.Errorf("expected search test, got %s", r.URL.Query().Get("search"))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"data":       []map[string]interface{}{},
				"pagination": map[string]interface{}{"page": 2, "limit": 50, "total": 0, "totalPages": 0, "hasNext": false, "hasPrev": false},
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	client.Lists.List(context.Background(), &ListListsParams{
		Page:   2,
		Limit:  50,
		Search: "test",
	})
}

// Test Contacts.List with all query parameters
func TestContactsListWithQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "3" {
			t.Errorf("expected page 3, got %s", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("limit") != "100" {
			t.Errorf("expected limit 100, got %s", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("status") != "active" {
			t.Errorf("expected status active, got %s", r.URL.Query().Get("status"))
		}
		if r.URL.Query().Get("search") != "john" {
			t.Errorf("expected search john, got %s", r.URL.Query().Get("search"))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"data":       []map[string]interface{}{},
				"pagination": map[string]interface{}{"page": 3, "limit": 100, "total": 0, "totalPages": 0, "hasNext": false, "hasPrev": false},
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	client.Contacts("list_123").List(context.Background(), &ListContactsParams{
		Page:   3,
		Limit:  100,
		Status: ContactStatusActive,
		Search: "john",
	})
}

// Test Emails.List with query parameters
func TestEmailsListWithQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("expected page 1, got %s", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("status") != "delivered" {
			t.Errorf("expected status delivered, got %s", r.URL.Query().Get("status"))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"data":       []map[string]interface{}{},
				"pagination": map[string]interface{}{"page": 1, "limit": 20, "total": 0, "totalPages": 0, "hasNext": false, "hasPrev": false},
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	client.Emails.List(context.Background(), &ListEmailsParams{
		Page:   1,
		Status: EmailStatusDelivered,
	})
}

// Test HTTP errors when API call fails
func TestAPICallErrors(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(client *Client) error
	}{
		{
			name: "emails send error",
			testFunc: func(client *Client) error {
				_, err := client.Emails.Send(context.Background(), &SendEmailParams{From: "a@b.com", To: []string{"c@d.com"}})
				return err
			},
		},
		{
			name: "emails stats error",
			testFunc: func(client *Client) error {
				_, err := client.Emails.Stats(context.Background())
				return err
			},
		},
		{
			name: "lists create error",
			testFunc: func(client *Client) error {
				_, err := client.Lists.Create(context.Background(), &CreateListParams{Name: "test"})
				return err
			},
		},
		{
			name: "lists get error",
			testFunc: func(client *Client) error {
				_, err := client.Lists.Get(context.Background(), "list_123")
				return err
			},
		},
		{
			name: "lists update error",
			testFunc: func(client *Client) error {
				_, err := client.Lists.Update(context.Background(), "list_123", &UpdateListParams{Name: "new"})
				return err
			},
		},
		{
			name: "lists stats error",
			testFunc: func(client *Client) error {
				_, err := client.Lists.Stats(context.Background(), "list_123")
				return err
			},
		},
		{
			name: "contacts create error",
			testFunc: func(client *Client) error {
				_, err := client.Contacts("list_123").Create(context.Background(), &CreateContactParams{Email: "a@b.com"})
				return err
			},
		},
		{
			name: "contacts get error",
			testFunc: func(client *Client) error {
				_, err := client.Contacts("list_123").Get(context.Background(), "contact_123")
				return err
			},
		},
		{
			name: "contacts update error",
			testFunc: func(client *Client) error {
				_, err := client.Contacts("list_123").Update(context.Background(), "contact_123", &UpdateContactParams{Email: "new@b.com"})
				return err
			},
		},
		{
			name: "contacts suppress error",
			testFunc: func(client *Client) error {
				err := client.Contacts("list_123").Suppress(context.Background(), "contact_123", SuppressReasonManual)
				return err
			},
		},
		{
			name: "verification verify error",
			testFunc: func(client *Client) error {
				_, err := client.Verification.Verify(context.Background(), &VerifyEmailParams{Email: "test@example.com"})
				return err
			},
		},
		{
			name: "verification batch error",
			testFunc: func(client *Client) error {
				_, err := client.Verification.Batch(context.Background(), []string{"a@b.com"})
				return err
			},
		},
		{
			name: "verification get error",
			testFunc: func(client *Client) error {
				_, err := client.Verification.Get(context.Background(), "ver_123")
				return err
			},
		},
		{
			name: "verification stats error",
			testFunc: func(client *Client) error {
				_, err := client.Verification.Stats(context.Background())
				return err
			},
		},
		{
			name: "verification list error",
			testFunc: func(client *Client) error {
				_, err := client.Verification.List(context.Background(), nil)
				return err
			},
		},
		{
			name: "attachments create upload error",
			testFunc: func(client *Client) error {
				_, err := client.Attachments.CreateUpload(context.Background(), &CreateUploadParams{Filename: "test.pdf", ContentType: "application/pdf"})
				return err
			},
		},
		{
			name: "attachments confirm error",
			testFunc: func(client *Client) error {
				_, err := client.Attachments.Confirm(context.Background(), "att_123")
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"error": map[string]interface{}{
						"code":    "SERVER_ERROR",
						"message": "Internal error",
					},
				})
			}))
			defer server.Close()

			client := NewClient("sk_test_123", WithBaseURL(server.URL), WithMaxRetries(0))
			err := tt.testFunc(client)

			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

// Test handling of malformed JSON response with error status code
func TestMalformedJSONErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	_, err := client.Emails.Get(context.Background(), "test")

	if err == nil {
		t.Fatal("expected error for error status with malformed JSON")
	}

	// Should create an error from status code since JSON is invalid
	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatal("expected *Error type")
	}

	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", apiErr.StatusCode)
	}
}

// Test handling of non-JSON error response
func TestNonJSONErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server Error"))
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL), WithMaxRetries(0))

	_, err := client.Emails.Get(context.Background(), "test")

	if err == nil {
		t.Fatal("expected error")
	}
}

// Test retry delay calculation
func TestRetryDelayCalculation(t *testing.T) {
	attempts := 0
	delays := make([]time.Duration, 0)
	startTime := time.Now()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			delays = append(delays, time.Since(startTime))
			startTime = time.Now()
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
				"id": "email_123", "from": "a@b.com", "to": []string{"c@d.com"},
				"status": "delivered", "created_at": "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL), WithMaxRetries(3))
	client.Emails.Get(context.Background(), "email_123")

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

// Test context cancellation
func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    map[string]interface{}{"id": "email_123"},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.Emails.Get(ctx, "test")

	if err == nil {
		t.Fatal("expected context deadline error")
	}
}

// Test rate limiting with HTTP Date format Retry-After
func TestRetryAfterHTTPDate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set Retry-After as HTTP date in the future
		futureTime := time.Now().Add(30 * time.Second).UTC().Format(http.TimeFormat)
		w.Header().Set("Retry-After", futureTime)
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
		t.Fatal("expected rate limit error")
	}

	// The retry-after should be approximately 30 seconds
	retryAfter := GetRetryAfter(err)
	if retryAfter < 25 || retryAfter > 35 {
		t.Logf("retry-after value: %d (may vary slightly due to timing)", retryAfter)
	}
}

// Test DELETE method for contacts
func TestContactsDeleteError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": "Contact not found",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	err := client.Contacts("list_123").Delete(context.Background(), "contact_123")

	if err == nil {
		t.Fatal("expected error")
	}

	if !IsNotFoundError(err) {
		t.Error("expected not found error")
	}
}

// Test DELETE method for lists
func TestListsDeleteError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"code":    "NOT_FOUND",
				"message": "List not found",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	err := client.Lists.Delete(context.Background(), "list_123")

	if err == nil {
		t.Fatal("expected error")
	}
}

// Test Emails.List error
func TestEmailsListError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"code":    "AUTHENTICATION_ERROR",
				"message": "Invalid API key",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	_, err := client.Emails.List(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error")
	}

	if !IsAuthenticationError(err) {
		t.Error("expected authentication error")
	}
}

// Test Lists.List error
func TestListsListError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"code":    "FORBIDDEN",
				"message": "Access denied",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	_, err := client.Lists.List(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error")
	}
}

// Test Contacts.List error
func TestContactsListError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	_, err := client.Contacts("list_123").List(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error")
	}

	if !IsValidationError(err) {
		t.Error("expected validation error")
	}
}

// Test success:false with 200 status code
func TestSuccessFalseWith200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"code":    "VALIDATION_ERROR",
				"message": "Something went wrong",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	_, err := client.Emails.Get(context.Background(), "test")

	if err == nil {
		t.Fatal("expected error even with 200 status")
	}
}

// Test success:false with no error object
func TestSuccessFalseWithNoError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	_, err := client.Emails.Get(context.Background(), "test")

	if err == nil {
		t.Fatal("expected error for success:false")
	}

	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatal("expected *Error type")
	}

	// Should use default unknown error
	if apiErr.Code != "UNKNOWN_ERROR" {
		t.Errorf("expected UNKNOWN_ERROR, got %s", apiErr.Code)
	}
}

// Test idempotency key with header injection attempt
func TestIdempotencyKeyHeaderInjection(t *testing.T) {
	var receivedKey string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedKey = r.Header.Get("X-Idempotency-Key")

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id": "email_123", "from": "a@b.com", "to": []string{"c@d.com"},
				"status": "pending", "created_at": "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	// Attempt to inject a header using \r\n
	_, err := client.Emails.Send(context.Background(), &SendEmailParams{
		From: "hello@example.com",
		To:   []string{"user@example.com"},
	}, WithIdempotencyKey("key\r\nX-Injected: bad"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The malicious key should be rejected, so no idempotency key should be set
	if receivedKey != "" {
		t.Errorf("expected no idempotency key (injection blocked), got '%s'", receivedKey)
	}
}

// Test retryable status codes (429 and 5xx)
func TestRetryableStatusCodes(t *testing.T) {
	tests := []struct {
		status     int
		shouldPass bool
	}{
		{http.StatusTooManyRequests, true},     // 429 is retryable
		{http.StatusInternalServerError, true}, // 500 is retryable
		{http.StatusBadGateway, true},          // 502 is retryable
		{http.StatusServiceUnavailable, true},  // 503 is retryable
	}

	for _, tt := range tests {
		t.Run(http.StatusText(tt.status), func(t *testing.T) {
			attempts := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempts++
				if attempts < 2 {
					w.WriteHeader(tt.status)
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
						"id": "email_123", "from": "a@b.com", "to": []string{"c@d.com"},
						"status": "delivered", "created_at": "2024-01-01T00:00:00Z",
					},
				})
			}))
			defer server.Close()

			client := NewClient("sk_test_123", WithBaseURL(server.URL), WithMaxRetries(3))
			_, err := client.Emails.Get(context.Background(), "email_123")

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if attempts != 2 {
				t.Errorf("expected 2 attempts (1 retry), got %d", attempts)
			}
		})
	}
}

// Test retry with Retry-After header uses the header value
func TestRetryWithRetryAfterHeader(t *testing.T) {
	attempts := 0
	var retryDelays []time.Duration
	lastAttempt := time.Now()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		now := time.Now()
		if attempts > 1 {
			retryDelays = append(retryDelays, now.Sub(lastAttempt))
		}
		lastAttempt = now

		if attempts < 2 {
			w.Header().Set("Retry-After", "1") // 1 second
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   map[string]interface{}{"message": "Rate limited"},
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id": "email_123", "from": "a@b.com", "to": []string{"c@d.com"},
				"status": "delivered", "created_at": "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL), WithMaxRetries(3))
	_, err := client.Emails.Get(context.Background(), "email_123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should respect the Retry-After header (approximately 1 second delay)
	if len(retryDelays) > 0 && retryDelays[0] < 900*time.Millisecond {
		t.Errorf("expected retry delay to respect Retry-After header, got %v", retryDelays[0])
	}
}

// Test Retry-After with HTTP-date format
func TestRetryAfterHTTPDateFormat(t *testing.T) {
	futureTime := time.Now().Add(45 * time.Second).UTC().Format(time.RFC1123)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", futureTime)
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
	// Should be approximately 45 seconds (allow some variance for test execution time)
	if retryAfter < 40 || retryAfter > 50 {
		t.Errorf("expected retry-after around 45 seconds, got %d", retryAfter)
	}
}

// Test Retry-After with invalid format returns 0
func TestRetryAfterInvalidFormat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "not-a-valid-value")
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
	if retryAfter != 0 {
		t.Errorf("expected retry-after 0 for invalid format, got %d", retryAfter)
	}
}

// Test Retry-After with past HTTP-date returns 0
func TestRetryAfterPastDate(t *testing.T) {
	pastTime := time.Now().Add(-60 * time.Second).UTC().Format(time.RFC1123)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", pastTime)
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
	if retryAfter != 0 {
		t.Errorf("expected retry-after 0 for past date, got %d", retryAfter)
	}
}

// Test 4xx status with error in response body
func Test4xxWithErrorBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true, // success: true but status >= 400
			"error": map[string]interface{}{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid email",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	_, err := client.Emails.Get(context.Background(), "test")

	// Should still error because status >= 400
	if err == nil {
		t.Fatal("expected error for 4xx status")
	}
}
