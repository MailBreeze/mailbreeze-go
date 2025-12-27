package mailbreeze

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("sk_test_123")

	if client == nil {
		t.Fatal("expected client to be created")
	}

	if client.Emails == nil {
		t.Error("expected Emails resource to be initialized")
	}
	if client.Lists == nil {
		t.Error("expected Lists resource to be initialized")
	}
	if client.Attachments == nil {
		t.Error("expected Attachments resource to be initialized")
	}
	if client.Verification == nil {
		t.Error("expected Verification resource to be initialized")
	}
}

func TestNewClientWithOptions(t *testing.T) {
	customClient := &http.Client{Timeout: 60 * time.Second}

	client := NewClient("sk_test_123",
		WithBaseURL("https://custom.api.com"),
		WithTimeout(60*time.Second),
		WithMaxRetries(5),
		WithHTTPClient(customClient),
	)

	if client == nil {
		t.Fatal("expected client to be created")
	}
}

func TestContacts(t *testing.T) {
	client := NewClient("sk_test_123")
	contacts := client.Contacts("list_123")

	if contacts == nil {
		t.Fatal("expected contacts resource to be created")
	}

	if contacts.listID != "list_123" {
		t.Errorf("expected listID to be 'list_123', got '%s'", contacts.listID)
	}
}

func TestEmailsSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/emails" {
			t.Errorf("expected /emails, got %s", r.URL.Path)
		}
		if r.Header.Get("X-API-Key") != "sk_test_123" {
			t.Errorf("expected X-API-Key header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json")
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":         "email_123",
				"from":       "hello@example.com",
				"to":         []string{"user@example.com"},
				"subject":    "Hello",
				"status":     "pending",
				"created_at": "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	email, err := client.Emails.Send(context.Background(), &SendEmailParams{
		From:    "hello@example.com",
		To:      []string{"user@example.com"},
		Subject: "Hello",
		HTML:    "<p>Hello</p>",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if email.ID != "email_123" {
		t.Errorf("expected ID 'email_123', got '%s'", email.ID)
	}
}

func TestEmailsSendWithIdempotencyKey(t *testing.T) {
	var receivedKey string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedKey = r.Header.Get("X-Idempotency-Key")

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":         "email_123",
				"from":       "hello@example.com",
				"to":         []string{"user@example.com"},
				"status":     "pending",
				"created_at": "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	_, err := client.Emails.Send(context.Background(), &SendEmailParams{
		From: "hello@example.com",
		To:   []string{"user@example.com"},
	}, WithIdempotencyKey("unique_key_123"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedKey != "unique_key_123" {
		t.Errorf("expected idempotency key 'unique_key_123', got '%s'", receivedKey)
	}
}

func TestEmailsList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/emails" {
			t.Errorf("expected /emails, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"id":         "email_1",
						"from":       "a@example.com",
						"to":         []string{"b@example.com"},
						"status":     "delivered",
						"created_at": "2024-01-01T00:00:00Z",
					},
				},
				"meta": map[string]interface{}{
					"page":        1,
					"limit":       20,
					"total":       1,
					"total_pages": 1,
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	result, err := client.Emails.List(context.Background(), &ListEmailsParams{
		Page:  1,
		Limit: 20,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(result.Items))
	}

	if result.Meta.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Meta.Total)
	}
}

func TestEmailsGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/emails/email_123" {
			t.Errorf("expected /emails/email_123, got %s", r.URL.Path)
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

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	email, err := client.Emails.Get(context.Background(), "email_123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if email.ID != "email_123" {
		t.Errorf("expected ID 'email_123', got '%s'", email.ID)
	}
}

func TestEmailsStats(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/emails/stats" {
			t.Errorf("expected /emails/stats, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"stats": map[string]interface{}{
					"total":         100,
					"sent":          95,
					"failed":        5,
					"transactional": 60,
					"marketing":     40,
					"successRate":   95.0,
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	stats, err := client.Emails.Stats(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.Total != 100 {
		t.Errorf("expected total 100, got %d", stats.Total)
	}
	if stats.Sent != 95 {
		t.Errorf("expected sent 95, got %d", stats.Sent)
	}
	if stats.SuccessRate != 95.0 {
		t.Errorf("expected successRate 95.0, got %f", stats.SuccessRate)
	}
}
