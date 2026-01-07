package nippou

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	domain "salesforce-mcp-server/internal/domain/nippou"
)

// ============================================================================
// Test Doubles
// ============================================================================

// MockRepository is a test double for domain.Repository.
type MockRepository struct {
	SaveFunc     func(n *domain.Nippou) error
	FindByIDFunc func(id domain.ID) (*domain.Nippou, error)
	FindByDateFunc func(date time.Time) ([]*domain.Nippou, error)
	DeleteFunc   func(id domain.ID) error
	SaveCalled   int
	LastSaved    *domain.Nippou
}

func (m *MockRepository) Save(n *domain.Nippou) error {
	m.SaveCalled++
	m.LastSaved = n
	if m.SaveFunc != nil {
		return m.SaveFunc(n)
	}
	return nil
}

func (m *MockRepository) FindByID(id domain.ID) (*domain.Nippou, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *MockRepository) FindByDate(date time.Time) ([]*domain.Nippou, error) {
	if m.FindByDateFunc != nil {
		return m.FindByDateFunc(date)
	}
	return nil, nil
}

func (m *MockRepository) Delete(id domain.ID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// ============================================================================
// UseCaseError Tests
// ============================================================================

func TestUseCaseError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *UseCaseError
		contains []string
	}{
		{
			name:     "error without cause",
			err:      &UseCaseError{Code: "CODE", Message: "message"},
			contains: []string{"CODE", "message"},
		},
		{
			name:     "error with cause",
			err:      &UseCaseError{Code: "CODE", Message: "message", Cause: errors.New("root cause")},
			contains: []string{"CODE", "message", "root cause"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := tt.err.Error()
			for _, s := range tt.contains {
				if !strings.Contains(errStr, s) {
					t.Errorf("Error() = %q, should contain %q", errStr, s)
				}
			}
		})
	}
}

func TestUseCaseError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := &UseCaseError{Code: "CODE", Message: "message", Cause: cause}

	unwrapped := errors.Unwrap(err)
	if unwrapped != cause {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, cause)
	}
}

func TestIsUseCaseError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"usecase error", &UseCaseError{Code: "CODE"}, true},
		{"wrapped usecase error", fmt.Errorf("wrapped: %w", &UseCaseError{Code: "CODE"}), true},
		{"standard error", errors.New("standard"), false},
		{"nil error", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUseCaseError(tt.err); got != tt.expected {
				t.Errorf("IsUseCaseError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsDomainViolation(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"domain violation", &UseCaseError{Code: ErrCodeDomainViolation}, true},
		{"invalid input", &UseCaseError{Code: ErrCodeInvalidInput}, false},
		{"repository error", &UseCaseError{Code: ErrCodeRepositoryError}, false},
		{"standard error", errors.New("standard"), false},
		{"nil error", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDomainViolation(tt.err); got != tt.expected {
				t.Errorf("IsDomainViolation() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// ============================================================================
// CreateInput Validation Tests
// ============================================================================

func TestCreateInput_Validate_Success(t *testing.T) {
	tests := []struct {
		name  string
		input *CreateInput
	}{
		{
			name: "minimal valid input",
			input: &CreateInput{
				Date:    "2026-01-08",
				Content: "Test content",
			},
		},
		{
			name: "full input with all optional fields",
			input: &CreateInput{
				Date:    "2026-01-08",
				Content: "Test content",
				Location: &LocationInput{
					Latitude:  35.6812,
					Longitude: 139.7671,
					Address:   "Tokyo Station",
				},
				Voice: &VoiceInput{
					Enabled:   true,
					ModelName: "gpt-4o-audio",
				},
				Tags: []string{"tag1", "tag2"},
			},
		},
		{
			name: "voice disabled with empty model",
			input: &CreateInput{
				Date:    "2026-01-08",
				Content: "Test content",
				Voice: &VoiceInput{
					Enabled:   false,
					ModelName: "",
				},
			},
		},
		{
			name: "boundary latitude values",
			input: &CreateInput{
				Date:    "2026-01-08",
				Content: "Test content",
				Location: &LocationInput{
					Latitude:  90.0,
					Longitude: 0,
				},
			},
		},
		{
			name: "boundary longitude values",
			input: &CreateInput{
				Date:    "2026-01-08",
				Content: "Test content",
				Location: &LocationInput{
					Latitude:  0,
					Longitude: 180.0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.input.Validate(); err != nil {
				t.Errorf("Validate() error = %v, want nil", err)
			}
		})
	}
}

func TestCreateInput_Validate_NilInput(t *testing.T) {
	var input *CreateInput
	err := input.Validate()
	if err != ErrNilInput {
		t.Errorf("Validate() error = %v, want ErrNilInput", err)
	}
}

func TestCreateInput_Validate_EmptyDate(t *testing.T) {
	tests := []string{"", "   ", "\t\n"}
	for _, date := range tests {
		input := &CreateInput{Date: date, Content: "content"}
		err := input.Validate()
		if err == nil {
			t.Errorf("Validate() with date=%q should return error", date)
		}
		if !strings.Contains(err.Error(), "date") {
			t.Errorf("Error should mention 'date': %v", err)
		}
	}
}

func TestCreateInput_Validate_EmptyContent(t *testing.T) {
	tests := []string{"", "   ", "\t\n"}
	for _, content := range tests {
		input := &CreateInput{Date: "2026-01-08", Content: content}
		err := input.Validate()
		if err == nil {
			t.Errorf("Validate() with content=%q should return error", content)
		}
		if !strings.Contains(err.Error(), "content") {
			t.Errorf("Error should mention 'content': %v", err)
		}
	}
}

func TestCreateInput_Validate_ContentTooLong(t *testing.T) {
	input := &CreateInput{
		Date:    "2026-01-08",
		Content: strings.Repeat("a", domain.MaxContentLength+1),
	}
	err := input.Validate()
	if err == nil {
		t.Error("Validate() with too long content should return error")
	}
	if !strings.Contains(err.Error(), "content") {
		t.Errorf("Error should mention 'content': %v", err)
	}
}

func TestCreateInput_Validate_TooManyTags(t *testing.T) {
	tags := make([]string, domain.MaxTagCount+1)
	for i := range tags {
		tags[i] = "tag"
	}
	input := &CreateInput{
		Date:    "2026-01-08",
		Content: "content",
		Tags:    tags,
	}
	err := input.Validate()
	if err == nil {
		t.Error("Validate() with too many tags should return error")
	}
	if !strings.Contains(err.Error(), "tags") {
		t.Errorf("Error should mention 'tags': %v", err)
	}
}

func TestCreateInput_Validate_TagTooLong(t *testing.T) {
	input := &CreateInput{
		Date:    "2026-01-08",
		Content: "content",
		Tags:    []string{"valid", strings.Repeat("a", domain.MaxTagLength+1)},
	}
	err := input.Validate()
	if err == nil {
		t.Error("Validate() with too long tag should return error")
	}
	if !strings.Contains(err.Error(), "tags") {
		t.Errorf("Error should mention 'tags': %v", err)
	}
}

func TestCreateInput_Validate_InvalidLatitude(t *testing.T) {
	tests := []float64{-91, 91, -100, 100, -90.001, 90.001}
	for _, lat := range tests {
		input := &CreateInput{
			Date:    "2026-01-08",
			Content: "content",
			Location: &LocationInput{
				Latitude:  lat,
				Longitude: 0,
			},
		}
		err := input.Validate()
		if err == nil {
			t.Errorf("Validate() with latitude=%f should return error", lat)
		}
		if !strings.Contains(err.Error(), "latitude") {
			t.Errorf("Error should mention 'latitude': %v", err)
		}
	}
}

func TestCreateInput_Validate_InvalidLongitude(t *testing.T) {
	tests := []float64{-181, 181, -200, 200, -180.001, 180.001}
	for _, lng := range tests {
		input := &CreateInput{
			Date:    "2026-01-08",
			Content: "content",
			Location: &LocationInput{
				Latitude:  0,
				Longitude: lng,
			},
		}
		err := input.Validate()
		if err == nil {
			t.Errorf("Validate() with longitude=%f should return error", lng)
		}
		if !strings.Contains(err.Error(), "longitude") {
			t.Errorf("Error should mention 'longitude': %v", err)
		}
	}
}

func TestCreateInput_Validate_AddressTooLong(t *testing.T) {
	input := &CreateInput{
		Date:    "2026-01-08",
		Content: "content",
		Location: &LocationInput{
			Latitude:  0,
			Longitude: 0,
			Address:   strings.Repeat("a", domain.MaxAddressLength+1),
		},
	}
	err := input.Validate()
	if err == nil {
		t.Error("Validate() with too long address should return error")
	}
	if !strings.Contains(err.Error(), "address") {
		t.Errorf("Error should mention 'address': %v", err)
	}
}

func TestCreateInput_Validate_VoiceEnabledEmptyModel(t *testing.T) {
	input := &CreateInput{
		Date:    "2026-01-08",
		Content: "content",
		Voice: &VoiceInput{
			Enabled:   true,
			ModelName: "",
		},
	}
	err := input.Validate()
	if err == nil {
		t.Error("Validate() with enabled voice but empty model should return error")
	}
	if !strings.Contains(err.Error(), "modelName") {
		t.Errorf("Error should mention 'modelName': %v", err)
	}
}

func TestCreateInput_Validate_ModelNameTooLong(t *testing.T) {
	input := &CreateInput{
		Date:    "2026-01-08",
		Content: "content",
		Voice: &VoiceInput{
			Enabled:   false,
			ModelName: strings.Repeat("a", domain.MaxModelNameLength+1),
		},
	}
	err := input.Validate()
	if err == nil {
		t.Error("Validate() with too long model name should return error")
	}
	if !strings.Contains(err.Error(), "modelName") {
		t.Errorf("Error should mention 'modelName': %v", err)
	}
}

// ============================================================================
// NewCreateUseCase Tests
// ============================================================================

func TestNewCreateUseCase_Success(t *testing.T) {
	repo := &MockRepository{}
	uc, err := NewCreateUseCase(repo)
	if err != nil {
		t.Fatalf("NewCreateUseCase() error = %v", err)
	}
	if uc == nil {
		t.Fatal("NewCreateUseCase() returned nil")
	}
}

func TestNewCreateUseCase_NilRepository(t *testing.T) {
	uc, err := NewCreateUseCase(nil)
	if err != ErrRepositoryNil {
		t.Errorf("NewCreateUseCase(nil) error = %v, want ErrRepositoryNil", err)
	}
	if uc != nil {
		t.Error("NewCreateUseCase(nil) should return nil UseCase")
	}
}

// ============================================================================
// Execute Tests - Success Cases
// ============================================================================

func TestExecute_MinimalInput_Success(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	input := &CreateInput{
		Date:    "2026-01-08",
		Content: "Test report content",
	}

	output, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if output == nil {
		t.Fatal("Execute() returned nil output")
	}
	if output.ID == "" {
		t.Error("Output.ID should not be empty")
	}
	if output.Date != "2026-01-08" {
		t.Errorf("Output.Date = %q, want '2026-01-08'", output.Date)
	}
	if output.Content != "Test report content" {
		t.Errorf("Output.Content = %q, want 'Test report content'", output.Content)
	}
	if output.Tags == nil {
		t.Error("Output.Tags should not be nil")
	}
	if len(output.Tags) != 0 {
		t.Errorf("Output.Tags should be empty, got %v", output.Tags)
	}
	if output.CreatedAt == "" {
		t.Error("Output.CreatedAt should not be empty")
	}
	if output.UpdatedAt == "" {
		t.Error("Output.UpdatedAt should not be empty")
	}
	if repo.SaveCalled != 1 {
		t.Errorf("Repository.Save() called %d times, want 1", repo.SaveCalled)
	}
}

func TestExecute_FullInput_Success(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	input := &CreateInput{
		Date:    "2026-01-08",
		Content: "Full report with all fields",
		Location: &LocationInput{
			Latitude:  35.6812,
			Longitude: 139.7671,
			Address:   "Tokyo Station",
		},
		Voice: &VoiceInput{
			Enabled:   true,
			ModelName: "gpt-4o-audio",
		},
		Tags: []string{"important", "sales"},
	}

	output, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Verify location
	if output.Location == nil {
		t.Fatal("Output.Location should not be nil")
	}
	if output.Location.Latitude != 35.6812 {
		t.Errorf("Output.Location.Latitude = %f, want 35.6812", output.Location.Latitude)
	}
	if output.Location.Longitude != 139.7671 {
		t.Errorf("Output.Location.Longitude = %f, want 139.7671", output.Location.Longitude)
	}
	if output.Location.Address != "Tokyo Station" {
		t.Errorf("Output.Location.Address = %q, want 'Tokyo Station'", output.Location.Address)
	}

	// Verify voice
	if output.Voice == nil {
		t.Fatal("Output.Voice should not be nil")
	}
	if !output.Voice.Enabled {
		t.Error("Output.Voice.Enabled should be true")
	}
	if output.Voice.ModelName != "gpt-4o-audio" {
		t.Errorf("Output.Voice.ModelName = %q, want 'gpt-4o-audio'", output.Voice.ModelName)
	}

	// Verify tags
	if len(output.Tags) != 2 {
		t.Errorf("Output.Tags count = %d, want 2", len(output.Tags))
	}
}

func TestExecute_ContentSanitization(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	input := &CreateInput{
		Date:    "2026-01-08",
		Content: "  Content with\x00null bytes and\r\nwindows line endings  ",
	}

	output, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Content should be sanitized by domain
	if strings.Contains(output.Content, "\x00") {
		t.Error("Output.Content should not contain null bytes")
	}
	if strings.Contains(output.Content, "\r") {
		t.Error("Output.Content should have normalized line endings")
	}
}

// ============================================================================
// Execute Tests - Error Cases
// ============================================================================

func TestExecute_NilContext(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	input := &CreateInput{Date: "2026-01-08", Content: "content"}
	_, err := uc.Execute(nil, input)
	if err != ErrContextNil {
		t.Errorf("Execute(nil, input) error = %v, want ErrContextNil", err)
	}
}

func TestExecute_NilInput(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	_, err := uc.Execute(context.Background(), nil)
	if err != ErrNilInput {
		t.Errorf("Execute(ctx, nil) error = %v, want ErrNilInput", err)
	}
}

func TestExecute_InvalidInputValidation(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	input := &CreateInput{Date: "", Content: "content"}
	_, err := uc.Execute(context.Background(), input)
	if err == nil {
		t.Error("Execute() with invalid input should return error")
	}
	if repo.SaveCalled != 0 {
		t.Error("Repository.Save() should not be called on input validation failure")
	}
}

func TestExecute_InvalidDateFormat(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	testCases := []string{
		"2026/01/08",
		"08-01-2026",
		"invalid",
		"2026-13-01",
	}
	for _, dateStr := range testCases {
		input := &CreateInput{Date: dateStr, Content: "content"}
		_, err := uc.Execute(context.Background(), input)
		if err == nil {
			t.Errorf("Execute() with date=%q should return error", dateStr)
		}
		if !IsDomainViolation(err) {
			t.Errorf("Execute() with date=%q should return domain violation, got: %v", dateStr, err)
		}
	}
}

func TestExecute_InvalidTagFormat(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	input := &CreateInput{
		Date:    "2026-01-08",
		Content: "content",
		Tags:    []string{"valid", "invalid tag with spaces"},
	}
	_, err := uc.Execute(context.Background(), input)
	if err == nil {
		t.Error("Execute() with invalid tag format should return error")
	}
	if !IsDomainViolation(err) {
		t.Errorf("Execute() should return domain violation, got: %v", err)
	}
}

func TestExecute_DuplicateTags(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	input := &CreateInput{
		Date:    "2026-01-08",
		Content: "content",
		Tags:    []string{"tag1", "TAG1"}, // Case-insensitive duplicate
	}
	_, err := uc.Execute(context.Background(), input)
	if err == nil {
		t.Error("Execute() with duplicate tags should return error")
	}
	if !IsDomainViolation(err) {
		t.Errorf("Execute() should return domain violation, got: %v", err)
	}
}

func TestExecute_RepositoryError(t *testing.T) {
	repoErr := errors.New("database connection failed")
	repo := &MockRepository{
		SaveFunc: func(n *domain.Nippou) error {
			return repoErr
		},
	}
	uc, _ := NewCreateUseCase(repo)

	input := &CreateInput{Date: "2026-01-08", Content: "content"}
	_, err := uc.Execute(context.Background(), input)
	if err == nil {
		t.Error("Execute() should return error when repository fails")
	}

	var ucErr *UseCaseError
	if !errors.As(err, &ucErr) {
		t.Fatalf("Error should be UseCaseError, got: %T", err)
	}
	if ucErr.Code != ErrCodeRepositoryError {
		t.Errorf("Error code = %q, want %q", ucErr.Code, ErrCodeRepositoryError)
	}
	if !errors.Is(err, repoErr) {
		t.Error("Error should wrap the original repository error")
	}
}

func TestExecute_ContextCancelled(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	input := &CreateInput{Date: "2026-01-08", Content: "content"}
	_, err := uc.Execute(ctx, input)
	if err == nil {
		t.Error("Execute() should return error when context is cancelled")
	}

	var ucErr *UseCaseError
	if !errors.As(err, &ucErr) {
		t.Fatalf("Error should be UseCaseError, got: %T", err)
	}
	if ucErr.Code != ErrCodeContextCancelled {
		t.Errorf("Error code = %q, want %q", ucErr.Code, ErrCodeContextCancelled)
	}
	if repo.SaveCalled != 0 {
		t.Error("Repository.Save() should not be called when context is cancelled")
	}
}

func TestExecute_ContextTimeout(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	// Let the timeout expire
	<-ctx.Done()

	input := &CreateInput{Date: "2026-01-08", Content: "content"}
	_, err := uc.Execute(ctx, input)
	if err == nil {
		t.Error("Execute() should return error when context times out")
	}
}

// ============================================================================
// Execute Tests - Domain Validation via Builder
// ============================================================================

func TestExecute_LocationDomainValidation(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	// This passes input validation but the domain location is created directly
	// Location input validation catches boundary issues at input level
	// but let's test that domain errors are properly propagated
	input := &CreateInput{
		Date:    "2026-01-08",
		Content: "content",
		Location: &LocationInput{
			Latitude:  45.0, // Valid
			Longitude: 90.0, // Valid
			Address:   "Valid address",
		},
	}

	output, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if output.Location == nil {
		t.Error("Output should have location")
	}
}

func TestExecute_VoiceDomainValidation(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	// Test disabled voice with model name
	input := &CreateInput{
		Date:    "2026-01-08",
		Content: "content",
		Voice: &VoiceInput{
			Enabled:   false,
			ModelName: "some-model",
		},
	}

	output, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if output.Voice == nil {
		t.Error("Output should have voice config")
	}
	if output.Voice.Enabled {
		t.Error("Voice should be disabled")
	}
}

// ============================================================================
// Execute Tests - Entity State Verification
// ============================================================================

func TestExecute_VerifySavedEntity(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	input := &CreateInput{
		Date:    "2026-01-08",
		Content: "Report content",
		Location: &LocationInput{
			Latitude:  35.6812,
			Longitude: 139.7671,
			Address:   "Tokyo",
		},
		Voice: &VoiceInput{
			Enabled:   true,
			ModelName: "model",
		},
		Tags: []string{"tag1", "tag2"},
	}

	_, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Verify the entity that was saved
	saved := repo.LastSaved
	if saved == nil {
		t.Fatal("No entity was saved")
	}
	if saved.ID().IsEmpty() {
		t.Error("Saved entity should have ID")
	}
	if saved.Content() != "Report content" {
		t.Errorf("Saved content = %q, want 'Report content'", saved.Content())
	}
	if saved.Location() == nil {
		t.Error("Saved entity should have location")
	}
	if saved.Voice() == nil {
		t.Error("Saved entity should have voice config")
	}
	if len(saved.Tags()) != 2 {
		t.Errorf("Saved entity should have 2 tags, got %d", len(saved.Tags()))
	}
}

// ============================================================================
// Output Mapping Tests
// ============================================================================

func TestMapToOutput_NilNippou(t *testing.T) {
	output := mapToOutput(nil)
	if output != nil {
		t.Error("mapToOutput(nil) should return nil")
	}
}

func TestMapToOutput_WithAllFields(t *testing.T) {
	// Create a nippou with all fields via builder
	loc, _ := domain.NewLocation(35.6812, 139.7671, "Tokyo")
	voice, _ := domain.NewVoiceConfig(true, "model")
	tag1, _ := domain.NewTag("tag1")
	tag2, _ := domain.NewTag("tag2")

	nippou, err := domain.NewNippouBuilder("2026-01-08", "content").
		WithLocation(loc).
		WithVoice(voice).
		WithTags([]domain.Tag{tag1, tag2}).
		Build()
	if err != nil {
		t.Fatalf("Failed to create test nippou: %v", err)
	}

	output := mapToOutput(nippou)

	if output.ID != nippou.ID().String() {
		t.Errorf("Output.ID = %q, want %q", output.ID, nippou.ID().String())
	}
	if output.Date != "2026-01-08" {
		t.Errorf("Output.Date = %q, want '2026-01-08'", output.Date)
	}
	if output.Content != "content" {
		t.Errorf("Output.Content = %q, want 'content'", output.Content)
	}
	if output.Location == nil {
		t.Error("Output.Location should not be nil")
	}
	if output.Voice == nil {
		t.Error("Output.Voice should not be nil")
	}
	if len(output.Tags) != 2 {
		t.Errorf("Output.Tags count = %d, want 2", len(output.Tags))
	}
}

func TestMapToOutput_WithoutOptionalFields(t *testing.T) {
	nippou, err := domain.NewNippou("2026-01-08", "content")
	if err != nil {
		t.Fatalf("Failed to create test nippou: %v", err)
	}

	output := mapToOutput(nippou)

	if output.Location != nil {
		t.Error("Output.Location should be nil")
	}
	if output.Voice != nil {
		t.Error("Output.Voice should be nil")
	}
	if output.Tags == nil {
		t.Error("Output.Tags should not be nil (should be empty slice)")
	}
	if len(output.Tags) != 0 {
		t.Errorf("Output.Tags should be empty, got %v", output.Tags)
	}
}

// ============================================================================
// Error Helper Function Tests
// ============================================================================

func TestNewInvalidInputError(t *testing.T) {
	err := NewInvalidInputError("field", "message")
	if err.Code != ErrCodeInvalidInput {
		t.Errorf("Code = %q, want %q", err.Code, ErrCodeInvalidInput)
	}
	if !strings.Contains(err.Message, "field") {
		t.Error("Message should contain field name")
	}
	if !strings.Contains(err.Message, "message") {
		t.Error("Message should contain message")
	}
}

func TestNewRepositoryError(t *testing.T) {
	cause := errors.New("db error")
	err := NewRepositoryError(cause)
	if err.Code != ErrCodeRepositoryError {
		t.Errorf("Code = %q, want %q", err.Code, ErrCodeRepositoryError)
	}
	if err.Cause != cause {
		t.Error("Cause should be preserved")
	}
}

func TestNewDomainViolationError(t *testing.T) {
	cause := domain.ErrEmptyContent
	err := NewDomainViolationError(cause)
	if err.Code != ErrCodeDomainViolation {
		t.Errorf("Code = %q, want %q", err.Code, ErrCodeDomainViolation)
	}
	if err.Cause != cause {
		t.Error("Cause should be preserved")
	}
}

// ============================================================================
// Concurrency Safety Tests
// ============================================================================

func TestExecute_ConcurrentCalls(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	const numCalls = 10
	results := make(chan error, numCalls)

	for i := 0; i < numCalls; i++ {
		go func(idx int) {
			input := &CreateInput{
				Date:    "2026-01-08",
				Content: "Concurrent content",
				Tags:    []string{"concurrent"},
			}
			_, err := uc.Execute(context.Background(), input)
			results <- err
		}(i)
	}

	for i := 0; i < numCalls; i++ {
		err := <-results
		if err != nil {
			t.Errorf("Concurrent Execute() error = %v", err)
		}
	}

	if repo.SaveCalled != numCalls {
		t.Errorf("Repository.Save() called %d times, want %d", repo.SaveCalled, numCalls)
	}
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestExecute_MaximumValidContent(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	input := &CreateInput{
		Date:    "2026-01-08",
		Content: strings.Repeat("a", domain.MaxContentLength),
	}

	output, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("Execute() with max content length error = %v", err)
	}
	if output == nil {
		t.Error("Output should not be nil")
	}
}

func TestExecute_MaximumValidTags(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	tags := make([]string, domain.MaxTagCount)
	for i := range tags {
		tags[i] = "tag" + string(rune('a'+i%26))
	}

	input := &CreateInput{
		Date:    "2026-01-08",
		Content: "content",
		Tags:    tags,
	}

	output, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("Execute() with max tags error = %v", err)
	}
	if len(output.Tags) != domain.MaxTagCount {
		t.Errorf("Output.Tags count = %d, want %d", len(output.Tags), domain.MaxTagCount)
	}
}

func TestExecute_BoundaryCoordinates(t *testing.T) {
	repo := &MockRepository{}
	uc, _ := NewCreateUseCase(repo)

	boundaryTests := []struct {
		name string
		lat  float64
		lng  float64
	}{
		{"south pole", -90.0, 0},
		{"north pole", 90.0, 0},
		{"date line west", 0, -180.0},
		{"date line east", 0, 180.0},
		{"all minimums", -90.0, -180.0},
		{"all maximums", 90.0, 180.0},
	}

	for _, bt := range boundaryTests {
		t.Run(bt.name, func(t *testing.T) {
			input := &CreateInput{
				Date:    "2026-01-08",
				Content: "content",
				Location: &LocationInput{
					Latitude:  bt.lat,
					Longitude: bt.lng,
				},
			}
			output, err := uc.Execute(context.Background(), input)
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if output.Location == nil {
				t.Error("Output.Location should not be nil")
			}
			if output.Location.Latitude != bt.lat {
				t.Errorf("Latitude = %f, want %f", output.Location.Latitude, bt.lat)
			}
			if output.Location.Longitude != bt.lng {
				t.Errorf("Longitude = %f, want %f", output.Location.Longitude, bt.lng)
			}
		})
	}
}
