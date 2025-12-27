package mailbreeze

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ContactsResource provides access to contact operations within a list.
type ContactsResource struct {
	client *HTTPClient
	listID string
}

// Create creates a new contact in the list.
func (r *ContactsResource) Create(ctx context.Context, params *CreateContactParams) (*Contact, error) {
	var contact Contact
	if err := r.client.Post(ctx, fmt.Sprintf("/api/v1/contact-lists/%s/contacts", r.listID), params, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

// List lists contacts in the list.
func (r *ContactsResource) List(ctx context.Context, params *ListContactsParams) (*ContactList, error) {
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
		if params.Search != "" {
			query.Set("search", params.Search)
		}
	}

	var result ContactList
	if err := r.client.Get(ctx, fmt.Sprintf("/api/v1/contact-lists/%s/contacts", r.listID), query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a contact by ID.
func (r *ContactsResource) Get(ctx context.Context, contactID string) (*Contact, error) {
	var contact Contact
	if err := r.client.Get(ctx, fmt.Sprintf("/api/v1/contact-lists/%s/contacts/%s", r.listID, contactID), nil, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

// Update updates a contact.
func (r *ContactsResource) Update(ctx context.Context, contactID string, params *UpdateContactParams) (*Contact, error) {
	var contact Contact
	if err := r.client.Put(ctx, fmt.Sprintf("/api/v1/contact-lists/%s/contacts/%s", r.listID, contactID), params, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

// Delete deletes a contact.
func (r *ContactsResource) Delete(ctx context.Context, contactID string) error {
	return r.client.Delete(ctx, fmt.Sprintf("/api/v1/contact-lists/%s/contacts/%s", r.listID, contactID))
}

// SuppressReason represents the reason for suppressing a contact.
type SuppressReason string

const (
	SuppressReasonManual       SuppressReason = "manual"
	SuppressReasonUnsubscribed SuppressReason = "unsubscribed"
	SuppressReasonBounced      SuppressReason = "bounced"
	SuppressReasonComplained   SuppressReason = "complained"
	SuppressReasonSpamTrap     SuppressReason = "spam_trap"
)

// Suppress suppresses a contact (adds to suppression list).
func (r *ContactsResource) Suppress(ctx context.Context, contactID string, reason SuppressReason) error {
	body := map[string]string{"reason": string(reason)}
	return r.client.Post(ctx, fmt.Sprintf("/api/v1/contact-lists/%s/contacts/%s/suppress", r.listID, contactID), body, nil)
}
