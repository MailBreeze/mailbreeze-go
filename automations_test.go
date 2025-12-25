package mailbreeze

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAutomationsEnroll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/automations/enroll" {
			t.Errorf("expected /automations/enroll, got %s", r.URL.Path)
		}

		var body EnrollParams
		json.NewDecoder(r.Body).Decode(&body)

		if body.AutomationID != "automation_123" {
			t.Errorf("expected automation_id 'automation_123', got '%s'", body.AutomationID)
		}

		if body.ContactID != "contact_123" {
			t.Errorf("expected contact_id 'contact_123', got '%s'", body.ContactID)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":            "enrollment_123",
				"automation_id": "automation_123",
				"contact_id":    "contact_123",
				"status":        "active",
				"current_step":  0,
				"created_at":    "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	enrollment, err := client.Automations.Enroll(context.Background(), &EnrollParams{
		AutomationID: "automation_123",
		ContactID:    "contact_123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if enrollment.ID != "enrollment_123" {
		t.Errorf("expected ID 'enrollment_123', got '%s'", enrollment.ID)
	}

	if enrollment.Status != EnrollmentStatusActive {
		t.Errorf("expected status 'active', got '%s'", enrollment.Status)
	}
}

func TestEnrollmentsList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/automations/enrollments" {
			t.Errorf("expected /automations/enrollments, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"id":            "enrollment_1",
						"automation_id": "automation_123",
						"contact_id":    "contact_1",
						"status":        "active",
						"current_step":  1,
						"created_at":    "2024-01-01T00:00:00Z",
					},
					{
						"id":            "enrollment_2",
						"automation_id": "automation_123",
						"contact_id":    "contact_2",
						"status":        "completed",
						"current_step":  5,
						"created_at":    "2024-01-02T00:00:00Z",
					},
				},
				"meta": map[string]interface{}{
					"page":        1,
					"limit":       20,
					"total":       2,
					"total_pages": 1,
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	result, err := client.Automations.Enrollments.List(context.Background(), &ListEnrollmentsParams{
		AutomationID: "automation_123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.Items))
	}

	if result.Meta.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Meta.Total)
	}
}

func TestEnrollmentsCancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/automations/enrollments/enrollment_123/cancel" {
			t.Errorf("expected /automations/enrollments/enrollment_123/cancel, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":        "enrollment_123",
				"cancelled": true,
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	result, err := client.Automations.Enrollments.Cancel(context.Background(), "enrollment_123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Cancelled {
		t.Error("expected cancelled to be true")
	}
}
