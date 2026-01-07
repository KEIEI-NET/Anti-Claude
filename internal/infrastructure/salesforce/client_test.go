package salesforce

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"salesforce-mcp-server/internal/domain/nippou"
)

// ============================================================================
// Mock HTTP Client
// ============================================================================

// MockHTTPClient implements HTTPDoer for testing.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// ============================================================================
// Mock Token Provider
// ============================================================================

// MockTokenProvider implements TokenProvider for testing.
type MockTokenProvider struct {
	Token string
	Err   error
}

func (m *MockTokenProvider) GetToken(ctx context.Context) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	return m.Token, nil
}

// ============================================================================
// Test Helpers
// ============================================================================

// newMockResponse creates a mock HTTP response with the given status and body.
func newMockResponse(status int, body interface{}) *http.Response {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewReader(bodyBytes)),
		Header:     make(http.Header),
	}
}

// newTestClient creates a Client with mock dependencies for testing.
func newTestClient(mockHTTP *MockHTTPClient) *Client {
	config := DefaultConfig("https://test.salesforce.com")
	config.MaxRetries = 0 // Disable retries for faster tests
	return NewClient(config, mockHTTP, &MockTokenProvider{Token: "test-token"})
}

// ============================================================================
// Client Tests
// ============================================================================

func TestClient_Get_Success(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify request
			if req.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", req.Method)
			}
			if req.Header.Get("Authorization") != "Bearer test-token" {
				t.Errorf("expected Bearer token, got %s", req.Header.Get("Authorization"))
			}

			return newMockResponse(200, map[string]string{"result": "ok"}), nil
		},
	}

	client := newTestClient(mockHTTP)
	var result map[string]string
	err := client.Get(context.Background(), "/test", &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["result"] != "ok" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestClient_Get_APIError(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return newMockResponse(404, []sfErrorResponse{
				{Message: "Record not found", ErrorCode: "NOT_FOUND"},
			}), nil
		},
	}

	client := newTestClient(mockHTTP)
	var result map[string]string
	err := client.Get(context.Background(), "/sobjects/Nippou__c/invalid", &result)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if !apiErr.IsNotFound() {
		t.Errorf("expected not found error, got status %d", apiErr.StatusCode)
	}
}

func TestClient_Post_Success(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if req.Header.Get("Content-Type") != "application/json" {
				t.Errorf("expected JSON content type")
			}

			return newMockResponse(201, CreateSObjectResult{
				ID:      "550e8400-e29b-41d4-a716-446655440003",
				Success: true,
			}), nil
		},
	}

	client := newTestClient(mockHTTP)
	var result CreateSObjectResult
	err := client.Post(context.Background(), "/sobjects/Nippou__c", map[string]string{"Name": "Test"}, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("expected success=true")
	}
	if result.ID != "550e8400-e29b-41d4-a716-446655440003" {
		t.Errorf("unexpected ID: %s", result.ID)
	}
}

func TestClient_TokenProviderError(t *testing.T) {
	config := DefaultConfig("https://test.salesforce.com")
	client := NewClient(config, &MockHTTPClient{}, &MockTokenProvider{
		Err: fmt.Errorf("token refresh failed"),
	})

	err := client.Get(context.Background(), "/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "failed to get auth token: token refresh failed" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestClient_Retry_RateLimited(t *testing.T) {
	callCount := 0
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			callCount++
			if callCount < 3 {
				return newMockResponse(429, []sfErrorResponse{
					{Message: "Rate limited", ErrorCode: "REQUEST_LIMIT_EXCEEDED"},
				}), nil
			}
			return newMockResponse(200, map[string]string{"ok": "true"}), nil
		},
	}

	config := DefaultConfig("https://test.salesforce.com")
	config.MaxRetries = 3
	config.RetryBaseDelay = 10 * time.Millisecond // Fast retries for testing
	client := NewClient(config, mockHTTP, &MockTokenProvider{Token: "test-token"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result map[string]string
	err := client.Get(ctx, "/test", &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

// ============================================================================
// Repository Tests
// ============================================================================

func TestNippouRepository_Save_Create(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// First call: Query to check if exists
			if req.Method == http.MethodGet {
				return newMockResponse(200, QueryResult{
					TotalSize: 0,
					Done:      true,
					Records:   []NippouSF{},
				}), nil
			}
			// Second call: Create
			if req.Method == http.MethodPost {
				return newMockResponse(201, CreateSObjectResult{
					ID:      "550e8400-e29b-41d4-a716-446655440003",
					Success: true,
				}), nil
			}
			return newMockResponse(400, nil), nil
		},
	}

	client := newTestClient(mockHTTP)
	repo := NewNippouRepository(client)

	// Create a valid Nippou
	n, err := nippou.NewNippou("2024-01-15", "Test content for today")
	if err != nil {
		t.Fatalf("failed to create Nippou: %v", err)
	}

	err = repo.Save(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNippouRepository_FindByDate(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return newMockResponse(200, QueryResult{
				TotalSize: 2,
				Done:      true,
				Records: []NippouSF{
					{
						ID:          "550e8400-e29b-41d4-a716-446655440001",
						Date:        "2024-01-15",
						Content:     "First entry",
						CreatedDate: "2024-01-15T10:00:00.000+0000",
					},
					{
						ID:          "550e8400-e29b-41d4-a716-446655440002",
						Date:        "2024-01-15",
						Content:     "Second entry",
						CreatedDate: "2024-01-15T14:00:00.000+0000",
					},
				},
			}), nil
		},
	}

	client := newTestClient(mockHTTP)
	repo := NewNippouRepository(client)

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	results, err := repo.FindByDate(date)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	if results[0].Content() != "First entry" {
		t.Errorf("unexpected content: %s", results[0].Content())
	}
}

func TestNippouRepository_Delete_Success(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return newMockResponse(204, nil), nil
		},
	}

	client := newTestClient(mockHTTP)
	repo := NewNippouRepository(client)

	id, _ := nippou.IDFromString("550e8400-e29b-41d4-a716-446655440000")
	err := repo.Delete(id)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNippouRepository_Delete_NotFound(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return newMockResponse(404, []sfErrorResponse{
				{Message: "Record not found", ErrorCode: "NOT_FOUND"},
			}), nil
		},
	}

	client := newTestClient(mockHTTP)
	repo := NewNippouRepository(client)

	id, _ := nippou.IDFromString("550e8400-e29b-41d4-a716-446655440000")
	err := repo.Delete(id)

	// Delete of non-existent record should not be an error
	if err != nil {
		t.Fatalf("delete of non-existent record should succeed, got: %v", err)
	}
}

// ============================================================================
// Model Mapping Tests
// ============================================================================

func TestNippouSF_ToDomain(t *testing.T) {
	sf := &NippouSF{
		ID:               "550e8400-e29b-41d4-a716-446655440000",
		Date:             "2024-01-15",
		Content:          "Test content",
		Latitude:         35.6762,
		Longitude:        139.6503,
		Address:          "Tokyo, Japan",
		VoiceOn:          true,
		VoiceModel:       "whisper-1",
		Tags:             "meeting,important",
		CreatedDate:      "2024-01-15T10:00:00.000+0000",
		LastModifiedDate: "2024-01-15T12:00:00.000+0000",
	}

	n, err := sf.ToDomain()
	if err != nil {
		t.Fatalf("failed to convert to domain: %v", err)
	}

	if n.Content() != "Test content" {
		t.Errorf("unexpected content: %s", n.Content())
	}
	if n.Location() == nil {
		t.Error("expected location to be set")
	}
	if n.Voice() == nil {
		t.Error("expected voice to be set")
	}
	if !n.Voice().Enabled() {
		t.Error("expected voice to be enabled")
	}
	if len(n.Tags()) != 2 {
		t.Errorf("expected 2 tags, got %d", len(n.Tags()))
	}
}

func TestFromDomain(t *testing.T) {
	builder := nippou.NewNippouBuilder("2024-01-15", "Test content")

	loc, _ := nippou.NewLocation(35.6762, 139.6503, "Tokyo")
	builder.WithLocation(loc)

	voice, _ := nippou.NewVoiceConfig(true, "whisper-1")
	builder.WithVoice(voice)

	n, err := builder.Build()
	if err != nil {
		t.Fatalf("failed to build Nippou: %v", err)
	}

	sf := FromDomain(n)

	if sf.Content != "Test content" {
		t.Errorf("unexpected content: %s", sf.Content)
	}
	if sf.Latitude != 35.6762 {
		t.Errorf("unexpected latitude: %f", sf.Latitude)
	}
	if !sf.VoiceOn {
		t.Error("expected VoiceOn to be true")
	}
}

func TestEscapeSOQL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"it's a test", "it''s a test"},
		{"back\\slash", "back\\\\slash"},
		{"mixed's\\test", "mixed''s\\\\test"},
	}

	for _, tc := range tests {
		result := EscapeSOQL(tc.input)
		if result != tc.expected {
			t.Errorf("EscapeSOQL(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

// ============================================================================
// Interface Compliance Test
// ============================================================================

func TestNippouRepository_ImplementsInterface(t *testing.T) {
	// This test ensures NippouRepository implements nippou.Repository
	var _ nippou.Repository = (*NippouRepository)(nil)
}
