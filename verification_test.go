package mailbreeze

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerificationVerify(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/email-verification/single" {
			t.Errorf("expected /api/v1/email-verification/single, got %s", r.URL.Path)
		}

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)

		if body["email"] != "test@example.com" {
			t.Errorf("expected email 'test@example.com', got '%s'", body["email"])
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"email":   "test@example.com",
				"result":  "valid",
				"isValid": true,
				"details": map[string]interface{}{
					"isDisposable":   false,
					"isRoleAccount":  false,
					"isFreeProvider": false,
					"hasMxRecords":   true,
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	result, err := client.Verification.Verify(context.Background(), &VerifyEmailParams{Email: "test@example.com"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.IsValid {
		t.Error("expected isValid to be true")
	}

	if result.Result != VerificationStatusValid {
		t.Errorf("expected result 'valid', got '%s'", result.Result)
	}
}

func TestVerificationBatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/email-verification/batch" {
			t.Errorf("expected /api/v1/email-verification/batch, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"verificationId":  "ver_123",
				"status":          "processing",
				"totalEmails":     3,
				"processedEmails": 0,
				"createdAt":       "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	result, err := client.Verification.Batch(context.Background(), []string{"a@example.com", "b@example.com", "c@example.com"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.VerificationID != "ver_123" {
		t.Errorf("expected verificationId 'ver_123', got '%s'", result.VerificationID)
	}

	if result.TotalEmails != 3 {
		t.Errorf("expected totalEmails 3, got %d", result.TotalEmails)
	}
}

func TestVerificationGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/email-verification/ver_123" {
			t.Errorf("expected /api/v1/email-verification/ver_123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"verificationId":  "ver_123",
				"status":          "completed",
				"totalEmails":     3,
				"processedEmails": 3,
				"createdAt":       "2024-01-01T00:00:00Z",
				"completedAt":     "2024-01-01T00:00:10Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	result, err := client.Verification.Get(context.Background(), "ver_123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "completed" {
		t.Errorf("expected status 'completed', got '%s'", result.Status)
	}

	if result.ProcessedEmails != 3 {
		t.Errorf("expected processedEmails 3, got %d", result.ProcessedEmails)
	}
}

func TestVerificationList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/email-verification" {
			t.Errorf("expected /api/v1/email-verification, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"verificationId":  "ver_1",
						"status":          "completed",
						"totalEmails":     100,
						"processedEmails": 100,
						"createdAt":       "2024-01-01T00:00:00Z",
					},
					{
						"verificationId":  "ver_2",
						"status":          "processing",
						"totalEmails":     50,
						"processedEmails": 25,
						"createdAt":       "2024-01-02T00:00:00Z",
					},
				},
				"pagination": map[string]interface{}{
					"page":       1,
					"limit":      20,
					"total":      2,
					"totalPages": 1,
					"hasNext":    false,
					"hasPrev":    false,
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	result, err := client.Verification.List(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Data) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.Data))
	}

	if result.Pagination.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Pagination.Total)
	}
}

func TestVerificationStats(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/email-verification/stats" {
			t.Errorf("expected /api/v1/email-verification/stats, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"totalVerified":      1000,
				"totalValid":         850,
				"totalInvalid":       100,
				"totalUnknown":       0,
				"totalVerifications": 50,
				"validPercentage":    85.0,
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	stats, err := client.Verification.Stats(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.TotalVerified != 1000 {
		t.Errorf("expected totalVerified 1000, got %d", stats.TotalVerified)
	}

	if stats.TotalValid != 850 {
		t.Errorf("expected totalValid 850, got %d", stats.TotalValid)
	}
}
