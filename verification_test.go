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
		if r.URL.Path != "/verification/verify" {
			t.Errorf("expected /verification/verify, got %s", r.URL.Path)
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
				"email":            "test@example.com",
				"status":           "valid",
				"is_valid":         true,
				"is_disposable":    false,
				"is_role_based":    false,
				"is_free_provider": false,
				"mx_found":         true,
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	result, err := client.Verification.Verify(context.Background(), "test@example.com")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.IsValid {
		t.Error("expected is_valid to be true")
	}

	if result.Status != VerificationStatusValid {
		t.Errorf("expected status 'valid', got '%s'", result.Status)
	}
}

func TestVerificationBatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/verification/batch" {
			t.Errorf("expected /verification/batch, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"verification_id": "ver_123",
				"status":          "processing",
				"total":           3,
				"processed":       0,
				"created_at":      "2024-01-01T00:00:00Z",
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
		t.Errorf("expected verification_id 'ver_123', got '%s'", result.VerificationID)
	}

	if result.Total != 3 {
		t.Errorf("expected total 3, got %d", result.Total)
	}
}

func TestVerificationGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/verification/ver_123" {
			t.Errorf("expected /verification/ver_123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"verification_id": "ver_123",
				"status":          "completed",
				"total":           3,
				"processed":       3,
				"created_at":      "2024-01-01T00:00:00Z",
				"completed_at":    "2024-01-01T00:00:10Z",
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

	if result.Processed != 3 {
		t.Errorf("expected processed 3, got %d", result.Processed)
	}
}

func TestVerificationStats(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/verification/stats" {
			t.Errorf("expected /verification/stats, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"total_verified": 1000,
				"valid_count":    850,
				"invalid_count":  100,
				"risky_count":    50,
				"unknown_count":  0,
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
		t.Errorf("expected total_verified 1000, got %d", stats.TotalVerified)
	}

	if stats.ValidCount != 850 {
		t.Errorf("expected valid_count 850, got %d", stats.ValidCount)
	}
}
