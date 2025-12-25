package mailbreeze

import (
	"context"
	"fmt"
)

// VerificationResource provides access to email verification operations.
type VerificationResource struct {
	client *HTTPClient
}

// Verify verifies a single email address.
func (r *VerificationResource) Verify(ctx context.Context, email string) (*VerificationResult, error) {
	var result VerificationResult
	body := map[string]string{"email": email}
	if err := r.client.Post(ctx, "/verification/verify", body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Batch starts a batch verification for multiple emails.
func (r *VerificationResource) Batch(ctx context.Context, emails []string) (*BatchVerificationResult, error) {
	var result BatchVerificationResult
	body := map[string][]string{"emails": emails}
	if err := r.client.Post(ctx, "/verification/batch", body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a batch verification status and results.
func (r *VerificationResource) Get(ctx context.Context, verificationID string) (*BatchVerificationResult, error) {
	var result BatchVerificationResult
	if err := r.client.Get(ctx, fmt.Sprintf("/verification/%s", verificationID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Stats returns verification statistics.
func (r *VerificationResource) Stats(ctx context.Context) (*VerificationStats, error) {
	var stats VerificationStats
	if err := r.client.Get(ctx, "/verification/stats", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
