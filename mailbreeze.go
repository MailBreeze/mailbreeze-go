// Package mailbreeze provides a Go client for the MailBreeze email platform API.
//
// Example usage:
//
//	client := mailbreeze.NewClient("sk_live_xxx")
//
//	email, err := client.Emails.Send(ctx, &mailbreeze.SendEmailParams{
//		From:    "hello@yourdomain.com",
//		To:      []string{"user@example.com"},
//		Subject: "Welcome!",
//		HTML:    "<h1>Welcome!</h1>",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(email.ID)
package mailbreeze

import (
	"net/http"
	"time"
)

// Version is the SDK version.
const Version = "1.0.3"

// DefaultBaseURL is the default API base URL.
const DefaultBaseURL = "https://api.mailbreeze.com/api/v1"

// DefaultTimeout is the default request timeout.
const DefaultTimeout = 30 * time.Second

// DefaultMaxRetries is the default number of retry attempts.
const DefaultMaxRetries = 3

// Client is the MailBreeze API client.
type Client struct {
	// Emails provides access to email operations.
	Emails *EmailsResource

	// Lists provides access to contact list operations.
	Lists *ListsResource

	// Attachments provides access to attachment operations.
	Attachments *AttachmentsResource

	// Verification provides access to email verification operations.
	Verification *VerificationResource

	httpClient *HTTPClient
}

// ClientOption is a function that configures the client.
type ClientOption func(*clientConfig)

type clientConfig struct {
	baseURL    string
	timeout    time.Duration
	maxRetries int
	httpClient *http.Client
}

// WithBaseURL sets a custom base URL.
func WithBaseURL(url string) ClientOption {
	return func(c *clientConfig) {
		c.baseURL = url
	}
}

// WithTimeout sets a custom request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *clientConfig) {
		c.timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retry attempts.
func WithMaxRetries(retries int) ClientOption {
	return func(c *clientConfig) {
		c.maxRetries = retries
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *clientConfig) {
		c.httpClient = client
	}
}

// NewClient creates a new MailBreeze API client.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	cfg := &clientConfig{
		baseURL:    DefaultBaseURL,
		timeout:    DefaultTimeout,
		maxRetries: DefaultMaxRetries,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.httpClient == nil {
		cfg.httpClient = &http.Client{
			Timeout: cfg.timeout,
		}
	}

	httpClient := newHTTPClient(apiKey, cfg.baseURL, cfg.maxRetries, cfg.httpClient)

	client := &Client{
		httpClient: httpClient,
	}

	// Initialize resources
	client.Emails = &EmailsResource{client: httpClient}
	client.Lists = &ListsResource{client: httpClient}
	client.Attachments = &AttachmentsResource{client: httpClient}
	client.Verification = &VerificationResource{client: httpClient}

	return client
}

// Contacts returns a ContactsResource scoped to the given list ID.
func (c *Client) Contacts(listID string) *ContactsResource {
	return &ContactsResource{
		client: c.httpClient,
		listID: listID,
	}
}

// String implements fmt.Stringer to prevent API key leakage in debug output.
func (c *Client) String() string {
	return "mailbreeze.Client{apiKey: [REDACTED]}"
}

// GoString implements fmt.GoStringer to prevent API key leakage in %#v output.
func (c *Client) GoString() string {
	return c.String()
}
