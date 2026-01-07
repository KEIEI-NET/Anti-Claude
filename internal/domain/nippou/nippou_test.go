package nippou

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// ============================================================================
// Test Helpers
// ============================================================================

// MockIDGenerator generates predictable IDs for testing.
type MockIDGenerator struct {
	ids   []string
	index int
}

func (m *MockIDGenerator) Generate() ID {
	if m.index >= len(m.ids) {
		return ID{value: "mock-id-overflow"}
	}
	id := ID{value: m.ids[m.index]}
	m.index++
	return id
}

func fixedTime() time.Time {
	return time.Date(2026, 1, 8, 12, 0, 0, 0, time.UTC)
}

// ============================================================================
// DomainError Tests
// ============================================================================

func TestDomainError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *DomainError
		expected string
	}{
		{
			name:     "with field",
			err:      &DomainError{Code: "CODE", Field: "field", Message: "message"},
			expected: "CODE: field - message",
		},
		{
			name:     "without field",
			err:      &DomainError{Code: "CODE", Message: "message"},
			expected: "CODE: message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"validation error", ErrEmptyContent, true},
		{"format error", ErrInvalidDateFormat, false},
		{"limit error", ErrContentTooLong, false},
		{"nil error", nil, false},
		{"standard error", errors.New("standard"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidationError(tt.err); got != tt.expected {
				t.Errorf("IsValidationError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// ============================================================================
// ID Tests
// ============================================================================

func TestNewID(t *testing.T) {
	id := NewID()
	if id.IsEmpty() {
		t.Error("NewID() should not return empty ID")
	}
	if id.String() == "" {
		t.Error("NewID().String() should not be empty")
	}
}

func TestIDFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid UUID", "550e8400-e29b-41d4-a716-446655440000", false},
		{"empty string", "", true},
		{"invalid format", "not-a-uuid", true},
		{"partial UUID", "550e8400-e29b", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := IDFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IDFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && id.String() != tt.input {
				t.Errorf("IDFromString() = %v, want %v", id.String(), tt.input)
			}
		})
	}
}

func TestID_Equals(t *testing.T) {
	id1, _ := IDFromString("550e8400-e29b-41d4-a716-446655440000")
	id2, _ := IDFromString("550e8400-e29b-41d4-a716-446655440000")
	id3, _ := IDFromString("660e8400-e29b-41d4-a716-446655440000")

	if !id1.Equals(id2) {
		t.Error("Equal IDs should be equal")
	}
	if id1.Equals(id3) {
		t.Error("Different IDs should not be equal")
	}
}

// ============================================================================
// Location Tests
// ============================================================================

func TestNewLocation_Success(t *testing.T) {
	loc, err := NewLocation(35.6812, 139.7671, "Tokyo Station")
	if err != nil {
		t.Fatalf("NewLocation() error = %v", err)
	}
	if loc.Latitude() != 35.6812 {
		t.Errorf("Latitude() = %v, want 35.6812", loc.Latitude())
	}
	if loc.Longitude() != 139.7671 {
		t.Errorf("Longitude() = %v, want 139.7671", loc.Longitude())
	}
	if loc.Address() != "Tokyo Station" {
		t.Errorf("Address() = %v, want 'Tokyo Station'", loc.Address())
	}
}

func TestNewLocation_InvalidLatitude(t *testing.T) {
	tests := []float64{-91, 91, -100, 100, -90.001, 90.001}
	for _, lat := range tests {
		_, err := NewLocation(lat, 0, "")
		if err != ErrInvalidLatitude {
			t.Errorf("NewLocation(%f, 0, '') error = %v, want ErrInvalidLatitude", lat, err)
		}
	}
}

func TestNewLocation_InvalidLongitude(t *testing.T) {
	tests := []float64{-181, 181, -200, 200, -180.001, 180.001}
	for _, lng := range tests {
		_, err := NewLocation(0, lng, "")
		if err != ErrInvalidLongitude {
			t.Errorf("NewLocation(0, %f, '') error = %v, want ErrInvalidLongitude", lng, err)
		}
	}
}

func TestNewLocation_BoundaryValues(t *testing.T) {
	tests := []struct {
		lat, lng float64
	}{
		{-90, -180},
		{90, 180},
		{0, 0},
		{-90, 0},
		{90, 0},
		{0, -180},
		{0, 180},
	}
	for _, tt := range tests {
		loc, err := NewLocation(tt.lat, tt.lng, "")
		if err != nil {
			t.Errorf("NewLocation(%f, %f, '') should succeed, got error: %v", tt.lat, tt.lng, err)
		}
		if loc == nil {
			t.Errorf("NewLocation(%f, %f, '') returned nil", tt.lat, tt.lng)
		}
	}
}

func TestNewLocation_AddressTooLong(t *testing.T) {
	longAddress := strings.Repeat("a", MaxAddressLength+1)
	_, err := NewLocation(0, 0, longAddress)
	if err != ErrAddressTooLong {
		t.Errorf("NewLocation() with long address error = %v, want ErrAddressTooLong", err)
	}
}

func TestNewLocation_SanitizesAddress(t *testing.T) {
	loc, err := NewLocation(0, 0, "  Address with\x00null  ")
	if err != nil {
		t.Fatalf("NewLocation() error = %v", err)
	}
	if strings.Contains(loc.Address(), "\x00") {
		t.Error("Address should not contain null bytes")
	}
	if loc.Address() != "Address withnull" {
		t.Errorf("Address() = %q, expected sanitized version", loc.Address())
	}
}

func TestLocation_NilSafe(t *testing.T) {
	var loc *Location
	if loc.Latitude() != 0 {
		t.Error("nil Location.Latitude() should return 0")
	}
	if loc.Longitude() != 0 {
		t.Error("nil Location.Longitude() should return 0")
	}
	if loc.Address() != "" {
		t.Error("nil Location.Address() should return empty string")
	}
}

func TestLocation_Equals(t *testing.T) {
	loc1, _ := NewLocation(35.6812, 139.7671, "Tokyo")
	loc2, _ := NewLocation(35.6812, 139.7671, "Tokyo")
	loc3, _ := NewLocation(35.6812, 139.7671, "Osaka")
	var nilLoc *Location

	if !loc1.Equals(loc2) {
		t.Error("Equal locations should be equal")
	}
	if loc1.Equals(loc3) {
		t.Error("Different locations should not be equal")
	}
	if loc1.Equals(nilLoc) {
		t.Error("Location should not equal nil")
	}
	if !nilLoc.Equals(nil) {
		t.Error("nil should equal nil")
	}
}

// ============================================================================
// VoiceConfig Tests
// ============================================================================

func TestNewVoiceConfig_Success(t *testing.T) {
	voice, err := NewVoiceConfig(true, "gpt-4o-audio")
	if err != nil {
		t.Fatalf("NewVoiceConfig() error = %v", err)
	}
	if !voice.Enabled() {
		t.Error("Enabled() should be true")
	}
	if voice.ModelName() != "gpt-4o-audio" {
		t.Errorf("ModelName() = %v, want 'gpt-4o-audio'", voice.ModelName())
	}
}

func TestNewVoiceConfig_DisabledWithEmptyModel(t *testing.T) {
	voice, err := NewVoiceConfig(false, "")
	if err != nil {
		t.Fatalf("NewVoiceConfig(false, '') error = %v", err)
	}
	if voice.Enabled() {
		t.Error("Enabled() should be false")
	}
}

func TestNewVoiceConfig_EnabledWithEmptyModel(t *testing.T) {
	_, err := NewVoiceConfig(true, "")
	if err != ErrEmptyModelName {
		t.Errorf("NewVoiceConfig(true, '') error = %v, want ErrEmptyModelName", err)
	}
}

func TestNewVoiceConfig_ModelNameTooLong(t *testing.T) {
	longModel := strings.Repeat("a", MaxModelNameLength+1)
	_, err := NewVoiceConfig(false, longModel)
	if err != ErrModelNameTooLong {
		t.Errorf("NewVoiceConfig() with long model error = %v, want ErrModelNameTooLong", err)
	}
}

func TestVoiceConfig_NilSafe(t *testing.T) {
	var voice *VoiceConfig
	if voice.Enabled() {
		t.Error("nil VoiceConfig.Enabled() should return false")
	}
	if voice.ModelName() != "" {
		t.Error("nil VoiceConfig.ModelName() should return empty string")
	}
}

func TestVoiceConfig_Equals(t *testing.T) {
	v1, _ := NewVoiceConfig(true, "model1")
	v2, _ := NewVoiceConfig(true, "model1")
	v3, _ := NewVoiceConfig(true, "model2")
	v4, _ := NewVoiceConfig(false, "model1")
	var nilVoice *VoiceConfig

	if !v1.Equals(v2) {
		t.Error("Equal VoiceConfigs should be equal")
	}
	if v1.Equals(v3) {
		t.Error("Different model names should not be equal")
	}
	if v1.Equals(v4) {
		t.Error("Different enabled states should not be equal")
	}
	if v1.Equals(nilVoice) {
		t.Error("VoiceConfig should not equal nil")
	}
	if !nilVoice.Equals(nil) {
		t.Error("nil should equal nil")
	}
}

// ============================================================================
// Tag Tests
// ============================================================================

func TestNewTag_Success(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"important", "important"},
		{"IMPORTANT", "important"},
		{"tag-with-dash", "tag-with-dash"},
		{"tag_with_underscore", "tag_with_underscore"},
		{"tag123", "tag123"},
		{"  trimmed  ", "trimmed"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tag, err := NewTag(tt.input)
			if err != nil {
				t.Errorf("NewTag(%q) error = %v", tt.input, err)
				return
			}
			if tag.String() != tt.expected {
				t.Errorf("NewTag(%q).String() = %q, want %q", tt.input, tag.String(), tt.expected)
			}
		})
	}
}

func TestNewTag_Empty(t *testing.T) {
	tests := []string{"", "   ", "\t\n"}
	for _, input := range tests {
		_, err := NewTag(input)
		if err != ErrEmptyTag {
			t.Errorf("NewTag(%q) error = %v, want ErrEmptyTag", input, err)
		}
	}
}

func TestNewTag_TooLong(t *testing.T) {
	longTag := strings.Repeat("a", MaxTagLength+1)
	_, err := NewTag(longTag)
	if err != ErrTagTooLong {
		t.Errorf("NewTag() with long tag error = %v, want ErrTagTooLong", err)
	}
}

func TestNewTag_InvalidFormat(t *testing.T) {
	tests := []string{
		"tag with spaces",
		"tag@special",
		"-starts-with-dash",
		"_starts-with-underscore",
		"tag!exclaim",
	}
	for _, input := range tests {
		_, err := NewTag(input)
		if err != ErrInvalidTagFormat {
			t.Errorf("NewTag(%q) error = %v, want ErrInvalidTagFormat", input, err)
		}
	}
}

func TestTag_Equals(t *testing.T) {
	tag1, _ := NewTag("important")
	tag2, _ := NewTag("IMPORTANT")
	tag3, _ := NewTag("sales")

	if !tag1.Equals(tag2) {
		t.Error("Tags should be case-insensitively equal")
	}
	if tag1.Equals(tag3) {
		t.Error("Different tags should not be equal")
	}
}

// ============================================================================
// Nippou Creation Tests
// ============================================================================

func TestNewNippou_Success(t *testing.T) {
	n, err := NewNippou("2026-01-08", "Today's report content")
	if err != nil {
		t.Fatalf("NewNippou() error = %v", err)
	}
	if n == nil {
		t.Fatal("NewNippou() returned nil")
	}
	if n.ID().IsEmpty() {
		t.Error("ID should not be empty")
	}
	if n.Content() != "Today's report content" {
		t.Errorf("Content() = %q, want 'Today's report content'", n.Content())
	}
	expectedDate, _ := time.Parse("2006-01-02", "2026-01-08")
	if !n.Date().Equal(expectedDate) {
		t.Errorf("Date() = %v, want %v", n.Date(), expectedDate)
	}
	if n.CreatedAt().IsZero() {
		t.Error("CreatedAt() should not be zero")
	}
	if n.UpdatedAt().IsZero() {
		t.Error("UpdatedAt() should not be zero")
	}
	if n.Tags() == nil {
		t.Error("Tags() should not be nil")
	}
	if len(n.Tags()) != 0 {
		t.Error("Tags() should be empty slice")
	}
}

func TestNewNippou_EmptyContent(t *testing.T) {
	_, err := NewNippou("2026-01-08", "")
	if err != ErrEmptyContent {
		t.Errorf("NewNippou() error = %v, want ErrEmptyContent", err)
	}
}

func TestNewNippou_WhitespaceOnlyContent(t *testing.T) {
	_, err := NewNippou("2026-01-08", "   \t\n  ")
	if err != ErrEmptyContent {
		t.Errorf("NewNippou() with whitespace content error = %v, want ErrEmptyContent", err)
	}
}

func TestNewNippou_ContentTooLong(t *testing.T) {
	longContent := strings.Repeat("a", MaxContentLength+1)
	_, err := NewNippou("2026-01-08", longContent)
	if err != ErrContentTooLong {
		t.Errorf("NewNippou() with long content error = %v, want ErrContentTooLong", err)
	}
}

func TestNewNippou_SanitizesContent(t *testing.T) {
	n, err := NewNippou("2026-01-08", "Content\x00with\rnull\r\nbytes")
	if err != nil {
		t.Fatalf("NewNippou() error = %v", err)
	}
	content := n.Content()
	if strings.Contains(content, "\x00") {
		t.Error("Content should not contain null bytes")
	}
	if strings.Contains(content, "\r") {
		t.Error("Content should have normalized line endings")
	}
}

func TestNewNippou_InvalidDateFormat(t *testing.T) {
	testCases := []string{
		"2026/01/08",
		"08-01-2026",
		"2026-1-8",
		"invalid",
		"",
		"2026-13-01",
		"2026-01-32",
	}
	for _, tc := range testCases {
		_, err := NewNippou(tc, "content")
		if err != ErrInvalidDateFormat {
			t.Errorf("NewNippou(%q, 'content') error = %v, want ErrInvalidDateFormat", tc, err)
		}
	}
}

// ============================================================================
// NippouBuilder Tests
// ============================================================================

func TestNippouBuilder_WithCustomIDGenerator(t *testing.T) {
	mockGen := &MockIDGenerator{ids: []string{"test-id-123"}}
	n, err := NewNippouBuilder("2026-01-08", "content").
		WithIDGenerator(mockGen).
		Build()

	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if n.ID().String() != "test-id-123" {
		t.Errorf("ID() = %q, want 'test-id-123'", n.ID().String())
	}
}

func TestNippouBuilder_WithCustomTimeFunc(t *testing.T) {
	n, err := NewNippouBuilder("2026-01-08", "content").
		WithTimeFunc(fixedTime).
		Build()

	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if !n.CreatedAt().Equal(fixedTime()) {
		t.Errorf("CreatedAt() = %v, want %v", n.CreatedAt(), fixedTime())
	}
}

func TestNippouBuilder_WithLocation(t *testing.T) {
	loc, _ := NewLocation(35.6812, 139.7671, "Tokyo")
	n, err := NewNippouBuilder("2026-01-08", "content").
		WithLocation(loc).
		Build()

	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if n.Location() == nil {
		t.Fatal("Location() should not be nil")
	}
	if n.Location().Latitude() != 35.6812 {
		t.Errorf("Location().Latitude() = %v, want 35.6812", n.Location().Latitude())
	}
}

func TestNippouBuilder_WithVoice(t *testing.T) {
	voice, _ := NewVoiceConfig(true, "gpt-4o")
	n, err := NewNippouBuilder("2026-01-08", "content").
		WithVoice(voice).
		Build()

	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if n.Voice() == nil {
		t.Fatal("Voice() should not be nil")
	}
	if !n.Voice().Enabled() {
		t.Error("Voice().Enabled() should be true")
	}
}

func TestNippouBuilder_WithTags(t *testing.T) {
	tag1, _ := NewTag("tag1")
	tag2, _ := NewTag("tag2")
	n, err := NewNippouBuilder("2026-01-08", "content").
		WithTags([]Tag{tag1, tag2}).
		Build()

	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if len(n.Tags()) != 2 {
		t.Errorf("Tags() count = %d, want 2", len(n.Tags()))
	}
}

func TestNippouBuilder_TooManyTags(t *testing.T) {
	tags := make([]Tag, MaxTagCount+1)
	for i := range tags {
		tags[i], _ = NewTag("tag" + string(rune('a'+i%26)))
	}
	_, err := NewNippouBuilder("2026-01-08", "content").
		WithTags(tags).
		Build()

	if err != ErrMaxTagsExceeded {
		t.Errorf("Build() with too many tags error = %v, want ErrMaxTagsExceeded", err)
	}
}

// ============================================================================
// Nippou Getter Tests (Nil Safety)
// ============================================================================

func TestNippou_NilSafeGetters(t *testing.T) {
	var n *Nippou

	if n.ID().String() != "" {
		t.Error("nil Nippou.ID().String() should return empty string")
	}
	if !n.Date().IsZero() {
		t.Error("nil Nippou.Date() should return zero time")
	}
	if n.Content() != "" {
		t.Error("nil Nippou.Content() should return empty string")
	}
	if n.Location() != nil {
		t.Error("nil Nippou.Location() should return nil")
	}
	if n.Voice() != nil {
		t.Error("nil Nippou.Voice() should return nil")
	}
	if n.Tags() != nil {
		t.Error("nil Nippou.Tags() should return nil")
	}
	if n.TagStrings() != nil {
		t.Error("nil Nippou.TagStrings() should return nil")
	}
	if !n.CreatedAt().IsZero() {
		t.Error("nil Nippou.CreatedAt() should return zero time")
	}
	if !n.UpdatedAt().IsZero() {
		t.Error("nil Nippou.UpdatedAt() should return zero time")
	}
	if n.TagCount() != 0 {
		t.Error("nil Nippou.TagCount() should return 0")
	}
	if n.HasTag("any") {
		t.Error("nil Nippou.HasTag() should return false")
	}
}

// ============================================================================
// Nippou Location Tests
// ============================================================================

func TestNippou_AttachLocation_Success(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	err := n.AttachLocation(35.6812, 139.7671, "Tokyo Station")
	if err != nil {
		t.Fatalf("AttachLocation() error = %v", err)
	}
	if n.Location() == nil {
		t.Fatal("Location() should not be nil after attach")
	}
	if n.Location().Latitude() != 35.6812 {
		t.Errorf("Location().Latitude() = %f, want 35.6812", n.Location().Latitude())
	}
}

func TestNippou_AttachLocation_InvalidLatitude(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	testCases := []float64{-91, 91, -100, 100}
	for _, lat := range testCases {
		err := n.AttachLocation(lat, 0, "")
		if err != ErrInvalidLatitude {
			t.Errorf("AttachLocation(%f, 0, '') error = %v, want ErrInvalidLatitude", lat, err)
		}
	}
}

func TestNippou_AttachLocation_InvalidLongitude(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	testCases := []float64{-181, 181, -200, 200}
	for _, lng := range testCases {
		err := n.AttachLocation(0, lng, "")
		if err != ErrInvalidLongitude {
			t.Errorf("AttachLocation(0, %f, '') error = %v, want ErrInvalidLongitude", lng, err)
		}
	}
}

func TestNippou_AttachLocation_BoundaryValues(t *testing.T) {
	boundaryTests := []struct {
		lat, lng float64
	}{
		{-90, -180},
		{90, 180},
		{0, 0},
	}
	for _, bt := range boundaryTests {
		n, _ := NewNippou("2026-01-08", "content")
		err := n.AttachLocation(bt.lat, bt.lng, "")
		if err != nil {
			t.Errorf("AttachLocation(%f, %f, '') error = %v", bt.lat, bt.lng, err)
		}
	}
}

func TestNippou_AttachLocation_NilReceiver(t *testing.T) {
	var n *Nippou
	err := n.AttachLocation(0, 0, "")
	if err != ErrNilNippou {
		t.Errorf("nil Nippou.AttachLocation() error = %v, want ErrNilNippou", err)
	}
}

func TestNippou_SetLocation(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	loc, _ := NewLocation(35.6812, 139.7671, "Tokyo")
	err := n.SetLocation(loc)
	if err != nil {
		t.Fatalf("SetLocation() error = %v", err)
	}
	if n.Location() == nil {
		t.Fatal("Location() should not be nil")
	}
}

func TestNippou_RemoveLocation(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	_ = n.AttachLocation(35.6812, 139.7671, "Tokyo")
	err := n.RemoveLocation()
	if err != nil {
		t.Fatalf("RemoveLocation() error = %v", err)
	}
	if n.Location() != nil {
		t.Error("Location() should be nil after removal")
	}
}

func TestNippou_Location_ReturnsDefensiveCopy(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	_ = n.AttachLocation(35.6812, 139.7671, "Tokyo")

	loc1 := n.Location()
	loc2 := n.Location()

	// They should be equal but not the same pointer
	if loc1 == loc2 {
		t.Error("Location() should return defensive copies, not same pointer")
	}
	if !loc1.Equals(loc2) {
		t.Error("Location() copies should be equal")
	}
}

// ============================================================================
// Nippou VoiceConfig Tests
// ============================================================================

func TestNippou_SetVoiceConfig(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	err := n.SetVoiceConfig(true, "gpt-4o-audio")
	if err != nil {
		t.Fatalf("SetVoiceConfig() error = %v", err)
	}
	if n.Voice() == nil {
		t.Fatal("Voice() should not be nil")
	}
	if !n.Voice().Enabled() {
		t.Error("Voice().Enabled() should be true")
	}
	if n.Voice().ModelName() != "gpt-4o-audio" {
		t.Errorf("Voice().ModelName() = %q, want 'gpt-4o-audio'", n.Voice().ModelName())
	}
}

func TestNippou_SetVoiceConfig_EnabledEmptyModel(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	err := n.SetVoiceConfig(true, "")
	if err != ErrEmptyModelName {
		t.Errorf("SetVoiceConfig(true, '') error = %v, want ErrEmptyModelName", err)
	}
}

func TestNippou_SetVoiceConfig_NilReceiver(t *testing.T) {
	var n *Nippou
	err := n.SetVoiceConfig(true, "model")
	if err != ErrNilNippou {
		t.Errorf("nil Nippou.SetVoiceConfig() error = %v, want ErrNilNippou", err)
	}
}

func TestNippou_SetVoice(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	voice, _ := NewVoiceConfig(true, "model")
	err := n.SetVoice(voice)
	if err != nil {
		t.Fatalf("SetVoice() error = %v", err)
	}
	if n.Voice() == nil {
		t.Fatal("Voice() should not be nil")
	}
}

func TestNippou_RemoveVoice(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	_ = n.SetVoiceConfig(true, "model")
	err := n.RemoveVoice()
	if err != nil {
		t.Fatalf("RemoveVoice() error = %v", err)
	}
	if n.Voice() != nil {
		t.Error("Voice() should be nil after removal")
	}
}

func TestNippou_Voice_ReturnsDefensiveCopy(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	_ = n.SetVoiceConfig(true, "model")

	v1 := n.Voice()
	v2 := n.Voice()

	if v1 == v2 {
		t.Error("Voice() should return defensive copies, not same pointer")
	}
	if !v1.Equals(v2) {
		t.Error("Voice() copies should be equal")
	}
}

// ============================================================================
// Nippou Tag Tests
// ============================================================================

func TestNippou_AddTag_Success(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	err := n.AddTag("important")
	if err != nil {
		t.Fatalf("AddTag() error = %v", err)
	}
	if len(n.Tags()) != 1 {
		t.Errorf("Tags() count = %d, want 1", len(n.Tags()))
	}
	if n.Tags()[0].String() != "important" {
		t.Errorf("Tags()[0] = %q, want 'important'", n.Tags()[0].String())
	}

	err = n.AddTag("sales")
	if err != nil {
		t.Fatalf("AddTag() error = %v", err)
	}
	if len(n.Tags()) != 2 {
		t.Errorf("Tags() count = %d, want 2", len(n.Tags()))
	}
}

func TestNippou_AddTag_Duplicate(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	_ = n.AddTag("important")
	err := n.AddTag("important")
	if err != ErrDuplicateTag {
		t.Errorf("AddTag() duplicate error = %v, want ErrDuplicateTag", err)
	}
	if len(n.Tags()) != 1 {
		t.Errorf("Tags() count = %d, want 1 after duplicate attempt", len(n.Tags()))
	}
}

func TestNippou_AddTag_CaseInsensitiveDuplicate(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	_ = n.AddTag("important")
	err := n.AddTag("IMPORTANT")
	if err != ErrDuplicateTag {
		t.Errorf("AddTag() case-insensitive duplicate error = %v, want ErrDuplicateTag", err)
	}
}

func TestNippou_AddTag_InvalidFormat(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	err := n.AddTag("tag with spaces")
	if err != ErrInvalidTagFormat {
		t.Errorf("AddTag() with invalid format error = %v, want ErrInvalidTagFormat", err)
	}
}

func TestNippou_AddTag_MaxExceeded(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	for i := 0; i < MaxTagCount; i++ {
		_ = n.AddTag("tag" + string(rune('a'+i)))
	}
	err := n.AddTag("overflow")
	if err != ErrMaxTagsExceeded {
		t.Errorf("AddTag() exceeding max error = %v, want ErrMaxTagsExceeded", err)
	}
}

func TestNippou_AddTag_NilReceiver(t *testing.T) {
	var n *Nippou
	err := n.AddTag("tag")
	if err != ErrNilNippou {
		t.Errorf("nil Nippou.AddTag() error = %v, want ErrNilNippou", err)
	}
}

func TestNippou_RemoveTag(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	_ = n.AddTag("important")
	_ = n.AddTag("sales")

	err := n.RemoveTag("important")
	if err != nil {
		t.Fatalf("RemoveTag() error = %v", err)
	}
	if len(n.Tags()) != 1 {
		t.Errorf("Tags() count = %d, want 1", len(n.Tags()))
	}
	if n.HasTag("important") {
		t.Error("HasTag('important') should be false after removal")
	}
	if !n.HasTag("sales") {
		t.Error("HasTag('sales') should still be true")
	}
}

func TestNippou_RemoveTag_NotFound(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	err := n.RemoveTag("nonexistent")
	if err != nil {
		t.Errorf("RemoveTag() for nonexistent tag error = %v, want nil", err)
	}
}

func TestNippou_HasTag(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	_ = n.AddTag("important")
	_ = n.AddTag("sales")

	if !n.HasTag("important") {
		t.Error("HasTag('important') should be true")
	}
	if !n.HasTag("IMPORTANT") {
		t.Error("HasTag('IMPORTANT') should be true (case-insensitive)")
	}
	if !n.HasTag("sales") {
		t.Error("HasTag('sales') should be true")
	}
	if n.HasTag("nonexistent") {
		t.Error("HasTag('nonexistent') should be false")
	}
}

func TestNippou_TagStrings(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	_ = n.AddTag("Important")
	_ = n.AddTag("Sales")

	tags := n.TagStrings()
	if len(tags) != 2 {
		t.Fatalf("TagStrings() count = %d, want 2", len(tags))
	}
	// Tags are normalized to lowercase
	if tags[0] != "important" || tags[1] != "sales" {
		t.Errorf("TagStrings() = %v, want ['important', 'sales']", tags)
	}
}

func TestNippou_Tags_ReturnsDefensiveCopy(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	_ = n.AddTag("tag1")

	tags1 := n.Tags()
	tags2 := n.Tags()

	// Modifying the returned slice should not affect the original
	if len(tags1) == 0 {
		t.Fatal("Tags() should not be empty")
	}
	tags1[0] = Tag{value: "modified"}
	if n.Tags()[0].String() == "modified" {
		t.Error("Tags() should return defensive copy")
	}
	if &tags1[0] == &tags2[0] {
		t.Error("Tags() should return different slice instances")
	}
}

// ============================================================================
// Nippou UpdateContent Tests
// ============================================================================

func TestNippou_UpdateContent_Success(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "original")
	originalUpdatedAt := n.UpdatedAt()

	time.Sleep(time.Millisecond)
	err := n.UpdateContent("updated content")
	if err != nil {
		t.Fatalf("UpdateContent() error = %v", err)
	}
	if n.Content() != "updated content" {
		t.Errorf("Content() = %q, want 'updated content'", n.Content())
	}
	if !n.UpdatedAt().After(originalUpdatedAt) {
		t.Error("UpdatedAt() should be updated")
	}
}

func TestNippou_UpdateContent_Empty(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "original")
	err := n.UpdateContent("")
	if err != ErrEmptyContent {
		t.Errorf("UpdateContent('') error = %v, want ErrEmptyContent", err)
	}
	if n.Content() != "original" {
		t.Error("Content() should not change on error")
	}
}

func TestNippou_UpdateContent_TooLong(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "original")
	err := n.UpdateContent(strings.Repeat("a", MaxContentLength+1))
	if err != ErrContentTooLong {
		t.Errorf("UpdateContent() with long content error = %v, want ErrContentTooLong", err)
	}
}

func TestNippou_UpdateContent_NilReceiver(t *testing.T) {
	var n *Nippou
	err := n.UpdateContent("content")
	if err != ErrNilNippou {
		t.Errorf("nil Nippou.UpdateContent() error = %v, want ErrNilNippou", err)
	}
}

// ============================================================================
// Reconstruct Tests
// ============================================================================

func TestReconstruct_Success(t *testing.T) {
	loc, _ := NewLocation(35.6812, 139.7671, "Tokyo")
	voice, _ := NewVoiceConfig(true, "model")
	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now()

	data := ReconstructedNippou{
		ID:        "550e8400-e29b-41d4-a716-446655440000",
		Date:      time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC),
		Content:   "Test content",
		Location:  loc,
		Voice:     voice,
		Tags:      []string{"tag1", "tag2"},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	n, err := Reconstruct(data)
	if err != nil {
		t.Fatalf("Reconstruct() error = %v", err)
	}
	if n.ID().String() != data.ID {
		t.Errorf("ID() = %q, want %q", n.ID().String(), data.ID)
	}
	if n.Content() != data.Content {
		t.Errorf("Content() = %q, want %q", n.Content(), data.Content)
	}
	if len(n.Tags()) != 2 {
		t.Errorf("Tags() count = %d, want 2", len(n.Tags()))
	}
}

func TestReconstruct_InvalidID(t *testing.T) {
	data := ReconstructedNippou{
		ID:      "invalid-id",
		Content: "content",
	}
	_, err := Reconstruct(data)
	if err == nil {
		t.Error("Reconstruct() with invalid ID should return error")
	}
}

func TestReconstruct_SkipsInvalidTags(t *testing.T) {
	data := ReconstructedNippou{
		ID:      "550e8400-e29b-41d4-a716-446655440000",
		Content: "content",
		Tags:    []string{"valid", "invalid tag with spaces", "another-valid"},
	}
	n, err := Reconstruct(data)
	if err != nil {
		t.Fatalf("Reconstruct() error = %v", err)
	}
	// Should have only the valid tags
	if len(n.Tags()) != 2 {
		t.Errorf("Tags() count = %d, want 2 (invalid skipped)", len(n.Tags()))
	}
}

// ============================================================================
// Helper Function Tests
// ============================================================================

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  trimmed  ", "trimmed"},
		{"null\x00byte", "nullbyte"},
		{"windows\r\nline", "windows\nline"},
		{"mac\rline", "mac\nline"},
		{"unix\nline", "unix\nline"},
		{"mixed\r\n\r\nlines", "mixed\n\nlines"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeString(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// UpdatedAt Tracking Tests
// ============================================================================

func TestNippou_UpdatedAt_TrackedOnMutations(t *testing.T) {
	n, _ := NewNippou("2026-01-08", "content")
	initial := n.UpdatedAt()

	// Test that each mutation updates the timestamp
	mutations := []struct {
		name string
		fn   func() error
	}{
		{"AttachLocation", func() error { return n.AttachLocation(0, 0, "") }},
		{"RemoveLocation", func() error { return n.RemoveLocation() }},
		{"SetVoiceConfig", func() error { return n.SetVoiceConfig(false, "") }},
		{"RemoveVoice", func() error { return n.RemoveVoice() }},
		{"AddTag", func() error { return n.AddTag("tag1") }},
		{"RemoveTag", func() error { return n.RemoveTag("tag1") }},
		{"UpdateContent", func() error { return n.UpdateContent("new content") }},
	}

	prev := initial
	for _, m := range mutations {
		time.Sleep(time.Millisecond)
		err := m.fn()
		if err != nil {
			t.Fatalf("%s() error = %v", m.name, err)
		}
		if !n.UpdatedAt().After(prev) {
			t.Errorf("%s() should update UpdatedAt", m.name)
		}
		prev = n.UpdatedAt()
	}
}
