package mailbreeze

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// VerificationResource provides access to email verification operations.
type VerificationResource struct {
	client *HTTPClient
}

// Verify verifies a single email address.
func (r *VerificationResource) Verify(ctx context.Context, params *VerifyEmailParams) (*VerificationResult, error) {
	var result VerificationResult
	if err := r.client.Post(ctx, "/api/v1/email-verification/single", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Batch starts a batch verification for multiple emails.
func (r *VerificationResource) Batch(ctx context.Context, emails []string) (*BatchVerificationResult, error) {
	var result BatchVerificationResult
	body := map[string][]string{"emails": emails}
	if err := r.client.Post(ctx, "/api/v1/email-verification/batch", body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a batch verification status and results.
func (r *VerificationResource) Get(ctx context.Context, verificationID string) (*BatchVerificationResult, error) {
	var result BatchVerificationResult
	if err := r.client.Get(ctx, fmt.Sprintf("/api/v1/email-verification/%s", verificationID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// List lists all batch verifications.
func (r *VerificationResource) List(ctx context.Context, params *ListVerificationsParams) (*VerificationsResponse, error) {
	query := url.Values{}

	if params != nil {
		if params.Page > 0 {
			query.Set("page", strconv.Itoa(params.Page))
		}
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Status != "" {
			query.Set("status", params.Status)
		}
	}

	var result VerificationsResponse
	if err := r.client.Get(ctx, "/api/v1/email-verification", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Stats returns verification statistics.
func (r *VerificationResource) Stats(ctx context.Context) (*VerificationStats, error) {
	var stats VerificationStats
	if err := r.client.Get(ctx, "/api/v1/email-verification/stats", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
