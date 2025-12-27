package mailbreeze

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// EmailsResource provides access to email operations.
type EmailsResource struct {
	client *HTTPClient
}

// Send sends an email.
func (r *EmailsResource) Send(ctx context.Context, params *SendEmailParams, opts ...RequestOption) (*Email, error) {
	var email Email
	if err := r.client.Post(ctx, "/emails", params, &email, opts...); err != nil {
		return nil, err
	}
	return &email, nil
}

// List lists emails with optional filtering.
func (r *EmailsResource) List(ctx context.Context, params *ListEmailsParams) (*EmailList, error) {
	query := url.Values{}

	if params != nil {
		if params.Status != "" {
			query.Set("status", string(params.Status))
		}
		if params.Page > 0 {
			query.Set("page", strconv.Itoa(params.Page))
		}
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
	}

	var result EmailList
	if err := r.client.Get(ctx, "/emails", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves an email by ID.
func (r *EmailsResource) Get(ctx context.Context, emailID string) (*Email, error) {
	var email Email
	if err := r.client.Get(ctx, fmt.Sprintf("/emails/%s", emailID), nil, &email); err != nil {
		return nil, err
	}
	return &email, nil
}

// Stats returns email statistics.
func (r *EmailsResource) Stats(ctx context.Context) (*EmailStats, error) {
	var response EmailStatsResponse
	if err := r.client.Get(ctx, "/emails/stats", nil, &response); err != nil {
		return nil, err
	}
	return &response.Stats, nil
}
