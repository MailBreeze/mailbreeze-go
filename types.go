package mailbreeze

import "time"

// PaginationMeta contains pagination information.
type PaginationMeta struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"totalPages"`
	HasNext    bool `json:"hasNext"`
	HasPrev    bool `json:"hasPrev"`
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
	CC          []string    `json:"cc,omitempty"`
	BCC         []string    `json:"bcc,omitempty"`
	Subject     string      `json:"subject,omitempty"`
	Status      EmailStatus `json:"status"`
	MessageID   string      `json:"messageId,omitempty"`
	TemplateID  string      `json:"templateId,omitempty"`
	CreatedAt   time.Time   `json:"createdAt"`
	SentAt      *time.Time  `json:"sentAt,omitempty"`
	DeliveredAt *time.Time  `json:"deliveredAt,omitempty"`
	OpenedAt    *time.Time  `json:"openedAt,omitempty"`
	ClickedAt   *time.Time  `json:"clickedAt,omitempty"`
}

// SendEmailResult is the result of sending an email.
type SendEmailResult struct {
	// MessageID is the unique message identifier returned by the API
	MessageID string `json:"messageId"`
}

// SendEmailParams are the parameters for sending an email.
type SendEmailParams struct {
	From          string            `json:"from"`
	To            []string          `json:"to"`
	Subject       string            `json:"subject,omitempty"`
	HTML          string            `json:"html,omitempty"`
	Text          string            `json:"text,omitempty"`
	TemplateID    string            `json:"templateId,omitempty"`
	Variables     map[string]any    `json:"variables,omitempty"`
	AttachmentIDs []string          `json:"attachmentIds,omitempty"`
	ReplyTo       string            `json:"replyTo,omitempty"`
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
	FromDate *time.Time  `json:"fromDate,omitempty"`
	ToDate   *time.Time  `json:"toDate,omitempty"`
}

// EmailList is a paginated list of emails.
type EmailList struct {
	Data       []Email        `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
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
	FirstName        string                 `json:"firstName,omitempty"`
	LastName         string                 `json:"lastName,omitempty"`
	PhoneNumber      string                 `json:"phoneNumber,omitempty"`
	Status           ContactStatus          `json:"status"`
	CustomFields     map[string]interface{} `json:"customFields,omitempty"`
	Source           string                 `json:"source,omitempty"`
	CreatedAt        time.Time              `json:"createdAt"`
	UpdatedAt        *time.Time             `json:"updatedAt,omitempty"`
	SubscribedAt     *time.Time             `json:"subscribedAt,omitempty"`
	UnsubscribedAt   *time.Time             `json:"unsubscribedAt,omitempty"`
	ConsentType      ConsentType            `json:"consentType,omitempty"`
	ConsentSource    string                 `json:"consentSource,omitempty"`
	ConsentTimestamp *time.Time             `json:"consentTimestamp,omitempty"`
	ConsentIpAddress string                 `json:"consentIpAddress,omitempty"`
}

// CreateContactParams are the parameters for creating a contact.
type CreateContactParams struct {
	Email            string                 `json:"email"`
	FirstName        string                 `json:"firstName,omitempty"`
	LastName         string                 `json:"lastName,omitempty"`
	PhoneNumber      string                 `json:"phoneNumber,omitempty"`
	CustomFields     map[string]interface{} `json:"customFields,omitempty"`
	Source           string                 `json:"source,omitempty"`
	ConsentType      ConsentType            `json:"consentType,omitempty"`
	ConsentSource    string                 `json:"consentSource,omitempty"`
	ConsentTimestamp *time.Time             `json:"consentTimestamp,omitempty"`
	ConsentIpAddress string                 `json:"consentIpAddress,omitempty"`
}

// UpdateContactParams are the parameters for updating a contact.
type UpdateContactParams struct {
	Email            string                 `json:"email,omitempty"`
	FirstName        string                 `json:"firstName,omitempty"`
	LastName         string                 `json:"lastName,omitempty"`
	PhoneNumber      string                 `json:"phoneNumber,omitempty"`
	CustomFields     map[string]interface{} `json:"customFields,omitempty"`
	ConsentType      ConsentType            `json:"consentType,omitempty"`
	ConsentSource    string                 `json:"consentSource,omitempty"`
	ConsentTimestamp *time.Time             `json:"consentTimestamp,omitempty"`
	ConsentIpAddress string                 `json:"consentIpAddress,omitempty"`
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
	Data       []Contact      `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// List represents a contact list.
type List struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description,omitempty"`
	ContactCount int        `json:"contactCount"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
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
	Data       []List         `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// ListStats contains list statistics.
type ListStats struct {
	TotalContacts        int `json:"totalContacts"`
	ActiveContacts       int `json:"activeContacts"`
	UnsubscribedContacts int `json:"unsubscribedContacts"`
	BouncedContacts      int `json:"bouncedContacts"`
	ComplainedContacts   int `json:"complainedContacts"`
	SuppressedContacts   int `json:"suppressedContacts"`
}

// VerificationStatus represents the verification result status.
type VerificationStatus string

const (
	VerificationStatusValid   VerificationStatus = "valid"
	VerificationStatusInvalid VerificationStatus = "invalid"
	VerificationStatusRisky   VerificationStatus = "risky"
	VerificationStatusUnknown VerificationStatus = "unknown"
)

// VerificationDetails contains additional verification details.
type VerificationDetails struct {
	IsFreeProvider bool `json:"isFreeProvider,omitempty"`
	IsDisposable   bool `json:"isDisposable,omitempty"`
	IsRoleAccount  bool `json:"isRoleAccount,omitempty"`
	HasMxRecords   bool `json:"hasMxRecords,omitempty"`
	IsSpamTrap     bool `json:"isSpamTrap,omitempty"`
}

// VerificationResult is the result of a single email verification.
type VerificationResult struct {
	Email     string               `json:"email"`
	IsValid   bool                 `json:"isValid"`
	Result    VerificationStatus   `json:"result"`
	Reason    string               `json:"reason,omitempty"`
	Cached    bool                 `json:"cached,omitempty"`
	RiskScore int                  `json:"riskScore,omitempty"`
	Details   *VerificationDetails `json:"details,omitempty"`
}

// BatchVerificationAnalytics contains analytics summary for batch verification.
type BatchVerificationAnalytics struct {
	Valid   int `json:"valid"`
	Invalid int `json:"invalid"`
	Risky   int `json:"risky"`
	Unknown int `json:"unknown"`
}

// BatchResults contains batch verification results grouped by category.
// This is returned when results are immediate (all cached).
type BatchResults struct {
	Clean   []string `json:"clean,omitempty"`
	Dirty   []string `json:"dirty,omitempty"`
	Unknown []string `json:"unknown,omitempty"`
}

// BatchVerificationResult is the result of a batch verification.
type BatchVerificationResult struct {
	VerificationID  string `json:"verificationId"`
	Status          string `json:"status"`
	TotalEmails     int    `json:"totalEmails"`
	ProcessedEmails int    `json:"processedEmails"`
	CreditsDeducted int    `json:"creditsDeducted"`
	// Results can be either []VerificationResult or BatchResults depending on API response
	Results     interface{}                 `json:"results,omitempty"`
	Analytics   *BatchVerificationAnalytics `json:"analytics,omitempty"`
	CreatedAt   time.Time                   `json:"createdAt"`
	CompletedAt *time.Time                  `json:"completedAt,omitempty"`
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

// ListVerificationsParams are the parameters for listing batch verifications.
type ListVerificationsParams struct {
	Page   int    `json:"page,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Status string `json:"status,omitempty"`
}

// VerificationsResponse is a paginated list of batch verifications.
type VerificationsResponse struct {
	Data       []BatchVerificationResult `json:"data"`
	Pagination PaginationMeta            `json:"pagination"`
}

// UploadURL contains the pre-signed upload URL.
type UploadURL struct {
	AttachmentID string    `json:"attachmentId"`
	UploadURL    string    `json:"uploadUrl"`
	UploadToken  string    `json:"uploadToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

// CreateUploadParams are the parameters for creating an upload URL.
type CreateUploadParams struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
	Inline      bool   `json:"inline,omitempty"`
}

// Attachment represents an attachment.
type Attachment struct {
	ID          string    `json:"id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"contentType"`
	Size        int64     `json:"size"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}
