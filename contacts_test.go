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
		if r.URL.Path != "/lists/list_123/contacts" {
			t.Errorf("expected /lists/list_123/contacts, got %s", r.URL.Path)
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
				"id":         "contact_123",
				"email":      "user@example.com",
				"status":     "subscribed",
				"created_at": "2024-01-01T00:00:00Z",
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
		if r.URL.Path != "/lists/list_123/contacts" {
			t.Errorf("expected /lists/list_123/contacts, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"id":         "contact_1",
						"email":      "user1@example.com",
						"status":     "subscribed",
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

	result, err := client.Contacts("list_123").List(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(result.Items))
	}
}

func TestContactsGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/lists/list_123/contacts/contact_123" {
			t.Errorf("expected /lists/list_123/contacts/contact_123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":         "contact_123",
				"email":      "user@example.com",
				"status":     "subscribed",
				"created_at": "2024-01-01T00:00:00Z",
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
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/lists/list_123/contacts/contact_123" {
			t.Errorf("expected /lists/list_123/contacts/contact_123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":         "contact_123",
				"email":      "updated@example.com",
				"status":     "subscribed",
				"created_at": "2024-01-01T00:00:00Z",
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
		if r.URL.Path != "/lists/list_123/contacts/contact_123" {
			t.Errorf("expected /lists/list_123/contacts/contact_123, got %s", r.URL.Path)
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
		if r.URL.Path != "/lists/list_123/contacts/contact_123/suppress" {
			t.Errorf("expected /lists/list_123/contacts/contact_123/suppress, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":         "contact_123",
				"email":      "user@example.com",
				"status":     "suppressed",
				"created_at": "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	contact, err := client.Contacts("list_123").Suppress(context.Background(), "contact_123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if contact.Status != ContactStatusSuppressed {
		t.Errorf("expected status 'suppressed', got '%s'", contact.Status)
	}
}
