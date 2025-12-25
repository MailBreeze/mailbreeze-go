package mailbreeze

import (
	"context"
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
	if err := r.client.Post(ctx, "/lists", params, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// List lists all contact lists.
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

	var result ListsResponse
	if err := r.client.Get(ctx, "/lists", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a contact list by ID.
func (r *ListsResource) Get(ctx context.Context, listID string) (*List, error) {
	var list List
	if err := r.client.Get(ctx, fmt.Sprintf("/lists/%s", listID), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// Update updates a contact list.
func (r *ListsResource) Update(ctx context.Context, listID string, params *UpdateListParams) (*List, error) {
	var list List
	if err := r.client.Patch(ctx, fmt.Sprintf("/lists/%s", listID), params, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// Delete deletes a contact list.
func (r *ListsResource) Delete(ctx context.Context, listID string) error {
	return r.client.Delete(ctx, fmt.Sprintf("/lists/%s", listID))
}

// Stats returns statistics for a contact list.
func (r *ListsResource) Stats(ctx context.Context, listID string) (*ListStats, error) {
	var stats ListStats
	if err := r.client.Get(ctx, fmt.Sprintf("/lists/%s/stats", listID), nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
