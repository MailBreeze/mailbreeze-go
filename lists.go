package mailbreeze

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// ListsResource provides access to contact list operations.
type ListsResource struct {
	client *HTTPClient
}

// Create creates a new contact list.
func (r *ListsResource) Create(ctx context.Context, params *CreateListParams) (*List, error) {
	var list List
	if err := r.client.Post(ctx, "/api/v1/contact-lists", params, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// List lists all contact lists.
// The API may return either an array or a paginated object, this method handles both.
func (r *ListsResource) List(ctx context.Context, params *ListListsParams) (*ListsResponse, error) {
	query := url.Values{}

	if params != nil {
		if params.Page > 0 {
			query.Set("page", strconv.Itoa(params.Page))
		}
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Search != "" {
			query.Set("search", params.Search)
		}
	}

	// Use json.RawMessage to handle polymorphic response
	var raw json.RawMessage
	if err := r.client.Get(ctx, "/api/v1/contact-lists", query, &raw); err != nil {
		return nil, err
	}

	// Try to unmarshal as array first
	var lists []List
	if err := json.Unmarshal(raw, &lists); err == nil {
		// Response was an array, create synthetic pagination
		return &ListsResponse{
			Data: lists,
			Pagination: PaginationMeta{
				Page:       1,
				Limit:      len(lists),
				Total:      len(lists),
				TotalPages: 1,
				HasNext:    false,
				HasPrev:    false,
			},
		}, nil
	}

	// Otherwise unmarshal as paginated object
	var result ListsResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal lists response: %w", err)
	}
	return &result, nil
}

// Get retrieves a contact list by ID.
func (r *ListsResource) Get(ctx context.Context, listID string) (*List, error) {
	var list List
	if err := r.client.Get(ctx, fmt.Sprintf("/api/v1/contact-lists/%s", listID), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// Update updates a contact list.
func (r *ListsResource) Update(ctx context.Context, listID string, params *UpdateListParams) (*List, error) {
	var list List
	if err := r.client.Put(ctx, fmt.Sprintf("/api/v1/contact-lists/%s", listID), params, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// Delete deletes a contact list.
func (r *ListsResource) Delete(ctx context.Context, listID string) error {
	return r.client.Delete(ctx, fmt.Sprintf("/api/v1/contact-lists/%s", listID))
}

// Stats returns statistics for a contact list.
func (r *ListsResource) Stats(ctx context.Context, listID string) (*ListStats, error) {
	var stats ListStats
	if err := r.client.Get(ctx, fmt.Sprintf("/api/v1/contact-lists/%s/stats", listID), nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
