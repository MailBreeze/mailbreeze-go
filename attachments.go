package mailbreeze

import "context"

// AttachmentsResource provides access to attachment operations.
type AttachmentsResource struct {
	client *HTTPClient
}

// CreateUpload creates a pre-signed upload URL.
func (r *AttachmentsResource) CreateUpload(ctx context.Context, params *CreateUploadParams) (*UploadURL, error) {
	var result UploadURL
	if err := r.client.Post(ctx, "/attachments/upload", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Confirm confirms an attachment upload.
func (r *AttachmentsResource) Confirm(ctx context.Context, attachmentID string) (*Attachment, error) {
	var attachment Attachment
	body := map[string]string{"attachment_id": attachmentID}
	if err := r.client.Post(ctx, "/attachments/confirm", body, &attachment); err != nil {
		return nil, err
	}
	return &attachment, nil
}
