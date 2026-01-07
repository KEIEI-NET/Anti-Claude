package salesforce

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ============================================================================
// TokenProvider Interface - Dependency Injection for Auth (DIP)
// ============================================================================

// TokenProvider abstracts the OAuth token retrieval mechanism.
// Implementations may cache tokens, refresh them, or fetch from external sources.
type TokenProvider interface {
	// GetToken returns a valid access token.
	// Implementations should handle token refresh internally.
	GetToken(ctx context.Context) (string, error)
}

// ============================================================================
// HTTP Client Interface - Testability (DIP)
// ============================================================================

// HTTPDoer abstracts http.Client for testability.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// ============================================================================
// Salesforce API Errors - Structured Error Handling
// ============================================================================

// APIError represents a Salesforce REST API error response.
type APIError struct {
	StatusCode int
	Message    string
	ErrorCode  string
	Fields     []string
}

func (e *APIError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("salesforce API error [%d] %s: %s", e.StatusCode, e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("salesforce API error [%d]: %s", e.StatusCode, e.Message)
}

// IsNotFound checks if the error is a 404 Not Found.
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsUnauthorized checks if the error is a 401 Unauthorized.
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsForbidden checks if the error is a 403 Forbidden.
func (e *APIError) IsForbidden() bool {
	return e.StatusCode == http.StatusForbidden
}

// IsRateLimited checks if the error is a 429 Too Many Requests.
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// sfErrorResponse represents the Salesforce error response JSON structure.
type sfErrorResponse struct {
	Message   string   `json:"message"`
	ErrorCode string   `json:"errorCode"`
	Fields    []string `json:"fields,omitempty"`
}

// ============================================================================
// Client Configuration
// ============================================================================

// ClientConfig holds configuration for the Salesforce client.
type ClientConfig struct {
	BaseURL        string        // e.g., "https://your-instance.salesforce.com"
	APIVersion     string        // e.g., "v59.0"
	Timeout        time.Duration // HTTP request timeout
	MaxRetries     int           // Maximum retry attempts for transient errors
	RetryBaseDelay time.Duration // Base delay for exponential backoff
}

// DefaultConfig returns sensible default configuration.
func DefaultConfig(baseURL string) *ClientConfig {
	return &ClientConfig{
		BaseURL:        baseURL,
		APIVersion:     "v59.0",
		Timeout:        30 * time.Second,
		MaxRetries:     3,
		RetryBaseDelay: 500 * time.Millisecond,
	}
}

// ============================================================================
// Salesforce REST Client
// ============================================================================

// Client provides low-level access to Salesforce REST API.
type Client struct {
	config        *ClientConfig
	httpClient    HTTPDoer
	tokenProvider TokenProvider
}

// NewClient creates a new Salesforce client with the given dependencies.
func NewClient(config *ClientConfig, httpClient HTTPDoer, tokenProvider TokenProvider) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: config.Timeout}
	}
	return &Client{
		config:        config,
		httpClient:    httpClient,
		tokenProvider: tokenProvider,
	}
}

// apiEndpoint constructs the full API endpoint URL.
func (c *Client) apiEndpoint(path string) string {
	return fmt.Sprintf("%s/services/data/%s%s", c.config.BaseURL, c.config.APIVersion, path)
}

// ============================================================================
// Core HTTP Methods
// ============================================================================

// Get performs an HTTP GET request.
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	return c.doRequest(ctx, http.MethodGet, path, nil, result)
}

// Post performs an HTTP POST request with JSON body.
func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
	return c.doRequest(ctx, http.MethodPost, path, body, result)
}

// Patch performs an HTTP PATCH request with JSON body.
func (c *Client) Patch(ctx context.Context, path string, body interface{}) error {
	return c.doRequest(ctx, http.MethodPatch, path, body, nil)
}

// Delete performs an HTTP DELETE request.
func (c *Client) Delete(ctx context.Context, path string) error {
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// Query executes a SOQL query and returns the result.
func (c *Client) Query(ctx context.Context, soql string, result interface{}) error {
	path := fmt.Sprintf("/query?q=%s", soql)
	return c.Get(ctx, path, result)
}

// ============================================================================
// Internal Request Execution
// ============================================================================

// doRequest executes an HTTP request with authentication and retry logic.
func (c *Client) doRequest(ctx context.Context, method, path string, body, result interface{}) error {
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := c.config.RetryBaseDelay * time.Duration(1<<uint(attempt-1))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		err := c.executeRequest(ctx, method, path, body, result)
		if err == nil {
			return nil
		}

		// Check if error is retryable
		if apiErr, ok := err.(*APIError); ok {
			if !c.isRetryable(apiErr) {
				return err
			}
			lastErr = err
			continue
		}

		// Context errors are not retryable
		if ctx.Err() != nil {
			return ctx.Err()
		}

		lastErr = err
	}

	return lastErr
}

// executeRequest performs a single HTTP request.
func (c *Client) executeRequest(ctx context.Context, method, path string, body, result interface{}) error {
	// Get auth token
	token, err := c.tokenProvider.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get auth token: %w", err)
	}

	// Build request
	url := c.apiEndpoint(path)
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		return c.parseAPIError(resp.StatusCode, respBody)
	}

	// Parse successful response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// parseAPIError converts HTTP error response to APIError.
func (c *Client) parseAPIError(statusCode int, body []byte) *APIError {
	apiErr := &APIError{StatusCode: statusCode}

	// Try to parse as Salesforce error response array
	var sfErrors []sfErrorResponse
	if err := json.Unmarshal(body, &sfErrors); err == nil && len(sfErrors) > 0 {
		apiErr.Message = sfErrors[0].Message
		apiErr.ErrorCode = sfErrors[0].ErrorCode
		apiErr.Fields = sfErrors[0].Fields
		return apiErr
	}

	// Try to parse as single error object
	var sfError sfErrorResponse
	if err := json.Unmarshal(body, &sfError); err == nil && sfError.Message != "" {
		apiErr.Message = sfError.Message
		apiErr.ErrorCode = sfError.ErrorCode
		apiErr.Fields = sfError.Fields
		return apiErr
	}

	// Fallback to raw body as message
	apiErr.Message = string(body)
	return apiErr
}

// isRetryable determines if an API error should trigger a retry.
func (c *Client) isRetryable(err *APIError) bool {
	switch err.StatusCode {
	case http.StatusTooManyRequests,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusBadGateway:
		return true
	default:
		return false
	}
}

// ============================================================================
// SObject Operations
// ============================================================================

// CreateSObjectResult represents the response from creating an SObject.
type CreateSObjectResult struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

// CreateSObject creates a new SObject record.
func (c *Client) CreateSObject(ctx context.Context, objectName string, record interface{}) (*CreateSObjectResult, error) {
	path := fmt.Sprintf("/sobjects/%s", objectName)
	var result CreateSObjectResult
	if err := c.Post(ctx, path, record, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSObject retrieves an SObject by ID.
func (c *Client) GetSObject(ctx context.Context, objectName, id string, result interface{}) error {
	path := fmt.Sprintf("/sobjects/%s/%s", objectName, id)
	return c.Get(ctx, path, result)
}

// UpdateSObject updates an existing SObject record.
func (c *Client) UpdateSObject(ctx context.Context, objectName, id string, record interface{}) error {
	path := fmt.Sprintf("/sobjects/%s/%s", objectName, id)
	return c.Patch(ctx, path, record)
}

// DeleteSObject deletes an SObject by ID.
func (c *Client) DeleteSObject(ctx context.Context, objectName, id string) error {
	path := fmt.Sprintf("/sobjects/%s/%s", objectName, id)
	return c.Delete(ctx, path)
}
