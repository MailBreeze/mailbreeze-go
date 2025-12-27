package mailbreeze

import "time"

// PaginationMeta contains pagination information.
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// EmailStatus represents the delivery status of an email.
type EmailStatus string

const (
	EmailStatusPending    EmailStatus = "pending"
	EmailStatusQueued     EmailStatus = "queued"
	EmailStatusSent       EmailStatus = "sent"
	EmailStatusDelivered  EmailStatus = "delivered"
	EmailStatusBounced    EmailStatus = "bounced"
	EmailStatusComplained EmailStatus = "complained"
	EmailStatusFailed     EmailStatus = "failed"
)

// Email represents an email object.
type Email struct {
	ID          string      `json:"id"`
	From        string      `json:"from"`
	To          []string    `json:"to"`
	Subject     string      `json:"subject,omitempty"`
	Status      EmailStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	SentAt      *time.Time  `json:"sent_at,omitempty"`
	DeliveredAt *time.Time  `json:"delivered_at,omitempty"`
}

// SendEmailParams are the parameters for sending an email.
type SendEmailParams struct {
	From          string            `json:"from"`
	To            []string          `json:"to"`
	Subject       string            `json:"subject,omitempty"`
	HTML          string            `json:"html,omitempty"`
	Text          string            `json:"text,omitempty"`
	TemplateID    string            `json:"template_id,omitempty"`
	Variables     map[string]any    `json:"variables,omitempty"`
	AttachmentIDs []string          `json:"attachment_ids,omitempty"`
	ReplyTo       string            `json:"reply_to,omitempty"`
	CC            []string          `json:"cc,omitempty"`
	BCC           []string          `json:"bcc,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
}

// ListEmailsParams are the parameters for listing emails.
type ListEmailsParams struct {
	Status   EmailStatus `json:"status,omitempty"`
	Page     int         `json:"page,omitempty"`
	Limit    int         `json:"limit,omitempty"`
	FromDate *time.Time  `json:"from_date,omitempty"`
	ToDate   *time.Time  `json:"to_date,omitempty"`
}

// EmailList is a paginated list of emails.
type EmailList struct {
	Items []Email        `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

// EmailStats contains email statistics.
type EmailStats struct {
	Total         int     `json:"total"`
	Sent          int     `json:"sent"`
	Failed        int     `json:"failed"`
	Transactional int     `json:"transactional"`
	Marketing     int     `json:"marketing"`
	SuccessRate   float64 `json:"successRate"`
}

// EmailStatsResponse is the wrapper for the stats API response.
type EmailStatsResponse struct {
	Stats EmailStats `json:"stats"`
}

// ContactStatus represents the subscription status of a contact.
type ContactStatus string

const (
	ContactStatusActive       ContactStatus = "active"
	ContactStatusUnsubscribed ContactStatus = "unsubscribed"
	ContactStatusBounced      ContactStatus = "bounced"
	ContactStatusComplained   ContactStatus = "complained"
	ContactStatusSuppressed   ContactStatus = "suppressed"
)

// ConsentType represents the type of consent obtained from a contact (NDPR compliance).
type ConsentType string

const (
	ConsentTypeExplicit           ConsentType = "explicit"
	ConsentTypeImplicit           ConsentType = "implicit"
	ConsentTypeLegitimateInterest ConsentType = "legitimate_interest"
)

// Contact represents a contact.
type Contact struct {
	ID               string                 `json:"id"`
	Email            string                 `json:"email"`
	FirstName        string                 `json:"first_name,omitempty"`
	LastName         string                 `json:"last_name,omitempty"`
	PhoneNumber      string                 `json:"phone_number,omitempty"`
	Status           ContactStatus          `json:"status"`
	CustomFields     map[string]interface{} `json:"custom_fields,omitempty"`
	Source           string                 `json:"source,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        *time.Time             `json:"updated_at,omitempty"`
	ConsentType      ConsentType            `json:"consent_type,omitempty"`
	ConsentSource    string                 `json:"consent_source,omitempty"`
	ConsentTimestamp *time.Time             `json:"consent_timestamp,omitempty"`
	ConsentIpAddress string                 `json:"consent_ip_address,omitempty"`
}

// CreateContactParams are the parameters for creating a contact.
type CreateContactParams struct {
	Email            string                 `json:"email"`
	FirstName        string                 `json:"first_name,omitempty"`
	LastName         string                 `json:"last_name,omitempty"`
	PhoneNumber      string                 `json:"phone_number,omitempty"`
	CustomFields     map[string]interface{} `json:"custom_fields,omitempty"`
	Source           string                 `json:"source,omitempty"`
	ConsentType      ConsentType            `json:"consent_type,omitempty"`
	ConsentSource    string                 `json:"consent_source,omitempty"`
	ConsentTimestamp *time.Time             `json:"consent_timestamp,omitempty"`
	ConsentIpAddress string                 `json:"consent_ip_address,omitempty"`
}

// UpdateContactParams are the parameters for updating a contact.
type UpdateContactParams struct {
	Email            string                 `json:"email,omitempty"`
	FirstName        string                 `json:"first_name,omitempty"`
	LastName         string                 `json:"last_name,omitempty"`
	PhoneNumber      string                 `json:"phone_number,omitempty"`
	CustomFields     map[string]interface{} `json:"custom_fields,omitempty"`
	ConsentType      ConsentType            `json:"consent_type,omitempty"`
	ConsentSource    string                 `json:"consent_source,omitempty"`
	ConsentTimestamp *time.Time             `json:"consent_timestamp,omitempty"`
	ConsentIpAddress string                 `json:"consent_ip_address,omitempty"`
}

// ListContactsParams are the parameters for listing contacts.
type ListContactsParams struct {
	Status ContactStatus `json:"status,omitempty"`
	Page   int           `json:"page,omitempty"`
	Limit  int           `json:"limit,omitempty"`
	Search string        `json:"search,omitempty"`
}

// ContactList is a paginated list of contacts.
type ContactList struct {
	Items []Contact      `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

// List represents a contact list.
type List struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description,omitempty"`
	ContactCount int        `json:"contact_count"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

// CreateListParams are the parameters for creating a list.
type CreateListParams struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateListParams are the parameters for updating a list.
type UpdateListParams struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ListListsParams are the parameters for listing lists.
type ListListsParams struct {
	Page   int    `json:"page,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Search string `json:"search,omitempty"`
}

// ListsResponse is a paginated list of lists.
type ListsResponse struct {
	Items []List         `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

// ListStats contains list statistics.
type ListStats struct {
	Total        int `json:"total"`
	Active       int `json:"active"`
	Unsubscribed int `json:"unsubscribed"`
	Bounced      int `json:"bounced"`
	Complained   int `json:"complained"`
}

// VerificationStatus represents the verification result status.
type VerificationStatus string

const (
	VerificationStatusValid   VerificationStatus = "valid"
	VerificationStatusInvalid VerificationStatus = "invalid"
	VerificationStatusRisky   VerificationStatus = "risky"
	VerificationStatusUnknown VerificationStatus = "unknown"
)

// VerificationResult is the result of a single email verification.
type VerificationResult struct {
	Email          string             `json:"email"`
	Status         VerificationStatus `json:"status"`
	IsValid        bool               `json:"is_valid"`
	IsDisposable   bool               `json:"is_disposable"`
	IsRoleBased    bool               `json:"is_role_based"`
	IsFreeProvider bool               `json:"is_free_provider"`
	MXFound        bool               `json:"mx_found"`
	SMTPCheck      *bool              `json:"smtp_check,omitempty"`
	Suggestion     string             `json:"suggestion,omitempty"`
}

// BatchVerificationResult is the result of a batch verification.
type BatchVerificationResult struct {
	VerificationID string               `json:"verification_id"`
	Status         string               `json:"status"`
	Total          int                  `json:"total"`
	Processed      int                  `json:"processed"`
	Results        []VerificationResult `json:"results,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
	CompletedAt    *time.Time           `json:"completed_at,omitempty"`
}

// VerificationStats contains verification statistics.
type VerificationStats struct {
	TotalVerified      int     `json:"totalVerified"`
	TotalValid         int     `json:"totalValid"`
	TotalInvalid       int     `json:"totalInvalid"`
	TotalUnknown       int     `json:"totalUnknown"`
	TotalVerifications int     `json:"totalVerifications"`
	ValidPercentage    float64 `json:"validPercentage"`
}

// VerifyEmailParams are the parameters for verifying a single email.
type VerifyEmailParams struct {
	Email string `json:"email"`
}

// UploadURL contains the pre-signed upload URL.
type UploadURL struct {
	AttachmentID string    `json:"attachment_id"`
	UploadURL    string    `json:"upload_url"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// CreateUploadParams are the parameters for creating an upload URL.
type CreateUploadParams struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

// Attachment represents an attachment.
type Attachment struct {
	ID          string    `json:"id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
