package mailbreeze

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAttachmentsCreateUpload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/attachments/presigned-url" {
			t.Errorf("expected /api/v1/attachments/presigned-url, got %s", r.URL.Path)
		}

		var body CreateUploadParams
		json.NewDecoder(r.Body).Decode(&body)

		if body.Filename != "document.pdf" {
			t.Errorf("expected filename 'document.pdf', got '%s'", body.Filename)
		}

		if body.ContentType != "application/pdf" {
			t.Errorf("expected contentType 'application/pdf', got '%s'", body.ContentType)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"attachmentId": "att_123",
				"uploadUrl":    "https://storage.example.com/upload/att_123",
				"uploadToken":  "token_123",
				"expiresAt":    "2024-01-01T01:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	result, err := client.Attachments.CreateUpload(context.Background(), &CreateUploadParams{
		Filename:    "document.pdf",
		ContentType: "application/pdf",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.AttachmentID != "att_123" {
		t.Errorf("expected attachmentId 'att_123', got '%s'", result.AttachmentID)
	}

	if result.UploadURL != "https://storage.example.com/upload/att_123" {
		t.Errorf("expected uploadUrl, got '%s'", result.UploadURL)
	}
}

func TestAttachmentsConfirm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/attachments/att_123/confirm" {
			t.Errorf("expected /api/v1/attachments/att_123/confirm, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":          "att_123",
				"filename":    "document.pdf",
				"contentType": "application/pdf",
				"size":        12345,
				"status":      "ready",
				"createdAt":   "2024-01-01T00:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("sk_test_123", WithBaseURL(server.URL))

	attachment, err := client.Attachments.Confirm(context.Background(), "att_123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if attachment.ID != "att_123" {
		t.Errorf("expected ID 'att_123', got '%s'", attachment.ID)
	}

	if attachment.Filename != "document.pdf" {
		t.Errorf("expected filename 'document.pdf', got '%s'", attachment.Filename)
	}

	if attachment.Status != "ready" {
		t.Errorf("expected status 'ready', got '%s'", attachment.Status)
	}
}
