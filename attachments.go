package mailbreeze

import (
	"context"
	"fmt"
)

// AttachmentsResource provides access to attachment operations.
type AttachmentsResource struct {
	client *HTTPClient
}

// CreateUpload creates a pre-signed upload URL.
func (r *AttachmentsResource) CreateUpload(ctx context.Context, params *CreateUploadParams) (*UploadURL, error) {
	var result UploadURL
	if err := r.client.Post(ctx, "/api/v1/attachments/presigned-url", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Confirm confirms an attachment upload.
func (r *AttachmentsResource) Confirm(ctx context.Context, attachmentID string) (*Attachment, error) {
	var attachment Attachment
	if err := r.client.Post(ctx, fmt.Sprintf("/api/v1/attachments/%s/confirm", attachmentID), nil, &attachment); err != nil {
		return nil, err
	}
	return &attachment, nil
}
