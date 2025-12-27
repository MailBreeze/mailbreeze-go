package mailbreeze

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListsCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/contact-lists" {
			t.Errorf("expected /api/v1/contact-lists, got %s", r.URL.Path)
		}

		var body CreateListParams
		json.NewDecoder(r.Body).Decode(&body)

		if body.Name != "My List" {
			t.Errorf("expected name 'My List', got '%s'", body.Name)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":        "list_123",
				"name":      "My List",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	list, err := client.Lists.Create(context.Background(), &CreateListParams{
		Name: "My List",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if list.ID != "list_123" {
		t.Errorf("expected ID 'list_123', got '%s'", list.ID)
	}
}

func TestListsList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/contact-lists" {
			t.Errorf("expected /api/v1/contact-lists, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":        "list_1",
						"name":      "List One",
						"createdAt": "2024-01-01T00:00:00Z",
					},
					{
						"id":        "list_2",
						"name":      "List Two",
						"createdAt": "2024-01-02T00:00:00Z",
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

	result, err := client.Lists.List(context.Background(), nil)

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

func TestListsListArrayResponse(t *testing.T) {
	// Test that we handle array response (no pagination wrapper)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": []map[string]interface{}{
				{
					"id":        "list_1",
					"name":      "List One",
					"createdAt": "2024-01-01T00:00:00Z",
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	result, err := client.Lists.List(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("expected 1 item, got %d", len(result.Data))
	}

	// Should have synthetic pagination
	if result.Pagination.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Pagination.Total)
	}
}

func TestListsGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/contact-lists/list_123" {
			t.Errorf("expected /api/v1/contact-lists/list_123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":        "list_123",
				"name":      "My List",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	list, err := client.Lists.Get(context.Background(), "list_123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if list.ID != "list_123" {
		t.Errorf("expected ID 'list_123', got '%s'", list.ID)
	}
}

func TestListsUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/contact-lists/list_123" {
			t.Errorf("expected /api/v1/contact-lists/list_123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":        "list_123",
				"name":      "Updated List",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	list, err := client.Lists.Update(context.Background(), "list_123", &UpdateListParams{
		Name: "Updated List",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if list.Name != "Updated List" {
		t.Errorf("expected name 'Updated List', got '%s'", list.Name)
	}
}

func TestListsDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/contact-lists/list_123" {
			t.Errorf("expected /api/v1/contact-lists/list_123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	err := client.Lists.Delete(context.Background(), "list_123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListsStats(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/contact-lists/list_123/stats" {
			t.Errorf("expected /api/v1/contact-lists/list_123/stats, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"totalContacts":        1000,
				"activeContacts":       900,
				"unsubscribedContacts": 50,
				"bouncedContacts":      25,
				"complainedContacts":   10,
				"suppressedContacts":   15,
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	stats, err := client.Lists.Stats(context.Background(), "list_123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.TotalContacts != 1000 {
		t.Errorf("expected totalContacts 1000, got %d", stats.TotalContacts)
	}
}
