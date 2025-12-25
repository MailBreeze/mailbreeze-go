package mailbreeze

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// HTTPClient handles HTTP requests to the MailBreeze API.
type HTTPClient struct {
	apiKey     string
	baseURL    string
	maxRetries int
	httpClient *http.Client
}

// apiResponse is the standard API response envelope.
type apiResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   *apiError       `json:"error,omitempty"`
	Meta    json.RawMessage `json:"meta,omitempty"`
}

// apiError is the error structure from the API.
type apiError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// requestOptions contains options for a single request.
type requestOptions struct {
	IdempotencyKey string
}

// RequestOption is a function that configures request options.
type RequestOption func(*requestOptions)

// WithIdempotencyKey sets the idempotency key for the request.
func WithIdempotencyKey(key string) RequestOption {
	return func(o *requestOptions) {
		o.IdempotencyKey = key
	}
}

func newHTTPClient(apiKey, baseURL string, maxRetries int, httpClient *http.Client) *HTTPClient {
	return &HTTPClient{
		apiKey:     apiKey,
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		maxRetries: maxRetries,
		httpClient: httpClient,
	}
}

// String implements fmt.Stringer to prevent API key leakage in debug output.
func (c *HTTPClient) String() string {
	return fmt.Sprintf("HTTPClient{baseURL: %q, maxRetries: %d, apiKey: [REDACTED]}", c.baseURL, c.maxRetries)
}

// GoString implements fmt.GoStringer to prevent API key leakage in %#v output.
func (c *HTTPClient) GoString() string {
	return c.String()
}

// Get performs a GET request.
func (c *HTTPClient) Get(ctx context.Context, path string, query url.Values, result interface{}) error {
	return c.request(ctx, http.MethodGet, path, query, nil, result, nil)
}

// Post performs a POST request.
func (c *HTTPClient) Post(ctx context.Context, path string, body, result interface{}, opts ...RequestOption) error {
	return c.request(ctx, http.MethodPost, path, nil, body, result, opts)
}

// Patch performs a PATCH request.
func (c *HTTPClient) Patch(ctx context.Context, path string, body, result interface{}) error {
	return c.request(ctx, http.MethodPatch, path, nil, body, result, nil)
}

// Delete performs a DELETE request.
func (c *HTTPClient) Delete(ctx context.Context, path string) error {
	return c.request(ctx, http.MethodDelete, path, nil, nil, nil, nil)
}

func (c *HTTPClient) request(
	ctx context.Context,
	method, path string,
	query url.Values,
	body, result interface{},
	opts []RequestOption,
) error {
	// Apply options
	reqOpts := &requestOptions{}
	for _, opt := range opts {
		opt(reqOpts)
	}

	// Build URL
	reqURL := c.baseURL + path
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}

	// Serialize body once (reused for retries)
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	var lastErr error
	maxAttempts := c.maxRetries + 1

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Create fresh body reader for each attempt
		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		// Create request
		req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		c.setHeaders(req, reqOpts)

		// Execute request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			if attempt < maxAttempts {
				time.Sleep(c.retryDelay(attempt, nil))
				continue
			}
			return lastErr
		}

		// Handle response
		apiErr, err := c.handleResponse(resp, result)
		if err != nil {
			return err
		}

		if apiErr != nil {
			lastErr = apiErr
			if !c.isRetryable(apiErr) || attempt >= maxAttempts {
				return apiErr
			}
			time.Sleep(c.retryDelay(attempt, apiErr))
			continue
		}

		return nil
	}

	return lastErr
}

func (c *HTTPClient) setHeaders(req *http.Request, opts *requestOptions) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("User-Agent", "mailbreeze-go/"+Version)

	if opts.IdempotencyKey != "" {
		// Validate idempotency key to prevent header injection
		if matched, _ := regexp.MatchString(`[\r\n]`, opts.IdempotencyKey); matched {
			// Skip setting invalid key - will be caught by validation elsewhere if needed
			return
		}
		req.Header.Set("X-Idempotency-Key", opts.IdempotencyKey)
	}
}

func (c *HTTPClient) handleResponse(resp *http.Response, result interface{}) (*Error, error) {
	defer resp.Body.Close()

	requestID := resp.Header.Get("X-Request-Id")
	retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))

	// Handle 204 No Content
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	// Read body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var apiResp apiResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		// Non-JSON response
		if resp.StatusCode >= 400 {
			return newErrorFromStatus(resp.StatusCode, "Unknown error", "", requestID, retryAfter), nil
		}
		return nil, nil
	}

	// Check for API error
	if !apiResp.Success || apiResp.Error != nil {
		statusCode := resp.StatusCode
		if resp.StatusCode < 400 {
			statusCode = http.StatusBadRequest // Default for success: false with 2xx
		}

		errMsg := "Unknown error"
		errCode := "UNKNOWN_ERROR"
		var details map[string]interface{}

		if apiResp.Error != nil {
			errMsg = apiResp.Error.Message
			errCode = apiResp.Error.Code
			details = apiResp.Error.Details
		}

		return newError(statusCode, errMsg, errCode, requestID, retryAfter, details), nil
	}

	// Check HTTP status
	if resp.StatusCode >= 400 {
		errMsg := "Unknown error"
		errCode := "UNKNOWN_ERROR"
		if apiResp.Error != nil {
			errMsg = apiResp.Error.Message
			errCode = apiResp.Error.Code
		}
		return newErrorFromStatus(resp.StatusCode, errMsg, errCode, requestID, retryAfter), nil
	}

	// Unmarshal data into result
	if result != nil && len(apiResp.Data) > 0 {
		if err := json.Unmarshal(apiResp.Data, result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response data: %w", err)
		}
	}

	return nil, nil
}

func (c *HTTPClient) isRetryable(err *Error) bool {
	if err == nil {
		return false
	}
	return err.StatusCode == http.StatusTooManyRequests || err.StatusCode >= 500
}

func (c *HTTPClient) retryDelay(attempt int, err *Error) time.Duration {
	// Use Retry-After header if available
	if err != nil && err.RetryAfter > 0 {
		return time.Duration(err.RetryAfter) * time.Second
	}

	// Exponential backoff: 1s, 2s, 4s...
	return time.Duration(1<<(attempt-1)) * time.Second
}

func parseRetryAfter(value string) int {
	if value == "" {
		return 0
	}

	// Try parsing as integer seconds first
	if seconds, err := strconv.Atoi(value); err == nil {
		return seconds
	}

	// Try parsing as HTTP-date (RFC 7231)
	if t, err := time.Parse(time.RFC1123, value); err == nil {
		delta := int(time.Until(t).Seconds())
		if delta > 0 {
			return delta
		}
		return 0
	}

	return 0
}
