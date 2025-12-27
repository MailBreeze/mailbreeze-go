package mailbreeze

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContactsCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/contact-lists/list_123/contacts" {
			t.Errorf("expected /api/v1/contact-lists/list_123/contacts, got %s", r.URL.Path)
		}

		var body CreateContactParams
		json.NewDecoder(r.Body).Decode(&body)

		if body.Email != "user@example.com" {
			t.Errorf("expected email 'user@example.com', got '%s'", body.Email)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":        "contact_123",
				"email":     "user@example.com",
				"status":    "active",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	contact, err := client.Contacts("list_123").Create(context.Background(), &CreateContactParams{
		Email: "user@example.com",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if contact.ID != "contact_123" {
		t.Errorf("expected ID 'contact_123', got '%s'", contact.ID)
	}
}

func TestContactsList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/contact-lists/list_123/contacts" {
			t.Errorf("expected /api/v1/contact-lists/list_123/contacts, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":        "contact_1",
						"email":     "user1@example.com",
						"status":    "active",
						"createdAt": "2024-01-01T00:00:00Z",
					},
				},
				"pagination": map[string]interface{}{
					"page":       1,
					"limit":      20,
					"total":      1,
					"totalPages": 1,
					"hasNext":    false,
					"hasPrev":    false,
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	result, err := client.Contacts("list_123").List(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("expected 1 item, got %d", len(result.Data))
	}
}

func TestContactsGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/contact-lists/list_123/contacts/contact_123" {
			t.Errorf("expected /api/v1/contact-lists/list_123/contacts/contact_123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":        "contact_123",
				"email":     "user@example.com",
				"status":    "active",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	contact, err := client.Contacts("list_123").Get(context.Background(), "contact_123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if contact.ID != "contact_123" {
		t.Errorf("expected ID 'contact_123', got '%s'", contact.ID)
	}
}

func TestContactsUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/contact-lists/list_123/contacts/contact_123" {
			t.Errorf("expected /api/v1/contact-lists/list_123/contacts/contact_123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":        "contact_123",
				"email":     "updated@example.com",
				"status":    "active",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	contact, err := client.Contacts("list_123").Update(context.Background(), "contact_123", &UpdateContactParams{
		Email: "updated@example.com",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if contact.Email != "updated@example.com" {
		t.Errorf("expected email 'updated@example.com', got '%s'", contact.Email)
	}
}

func TestContactsDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/contact-lists/list_123/contacts/contact_123" {
			t.Errorf("expected /api/v1/contact-lists/list_123/contacts/contact_123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	err := client.Contacts("list_123").Delete(context.Background(), "contact_123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContactsSuppress(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/contact-lists/list_123/contacts/contact_123/suppress" {
			t.Errorf("expected /api/v1/contact-lists/list_123/contacts/contact_123/suppress, got %s", r.URL.Path)
		}

		// Verify reason is in body
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["reason"] != "manual" {
			t.Errorf("expected reason 'manual', got '%s'", body["reason"])
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":        "contact_123",
				"email":     "user@example.com",
				"status":    "suppressed",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	err := client.Contacts("list_123").Suppress(context.Background(), "contact_123", SuppressReasonManual)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
