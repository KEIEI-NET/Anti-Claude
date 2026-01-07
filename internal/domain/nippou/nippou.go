package nippou

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// ============================================================================
// Constants - Security & Validation Limits
// ============================================================================

const (
	// MaxContentLength is the maximum allowed length for content (64KB).
	MaxContentLength = 65536
	// MaxAddressLength is the maximum allowed length for address strings.
	MaxAddressLength = 500
	// MaxModelNameLength is the maximum allowed length for model names.
	MaxModelNameLength = 100
	// MaxTagLength is the maximum allowed length for a single tag.
	MaxTagLength = 50
	// MaxTagCount is the maximum number of tags allowed per Nippou.
	MaxTagCount = 20
	// MinLatitude is the minimum valid latitude.
	MinLatitude = -90.0
	// MaxLatitude is the maximum valid latitude.
	MaxLatitude = 90.0
	// MinLongitude is the minimum valid longitude.
	MinLongitude = -180.0
	// MaxLongitude is the maximum valid longitude.
	MaxLongitude = 180.0
)

// ============================================================================
// Domain Errors - Structured Error Handling
// ============================================================================

// DomainError represents a domain-specific error with code and context.
type DomainError struct {
	Code    string
	Message string
	Field   string
}

func (e *DomainError) Error() string {
	if e.Field != "" {
		return e.Code + ": " + e.Field + " - " + e.Message
	}
	return e.Code + ": " + e.Message
}

// Error codes for domain errors.
const (
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeInvalidFormat = "INVALID_FORMAT"
	ErrCodeDuplicate     = "DUPLICATE_ERROR"
	ErrCodeLimitExceeded = "LIMIT_EXCEEDED"
	ErrCodeNilReceiver   = "NIL_RECEIVER"
)

// Predefined domain errors for common validation failures.
var (
	ErrEmptyContent       = &DomainError{Code: ErrCodeValidation, Field: "content", Message: "content cannot be empty"}
	ErrContentTooLong     = &DomainError{Code: ErrCodeLimitExceeded, Field: "content", Message: "content exceeds maximum length"}
	ErrInvalidDateFormat  = &DomainError{Code: ErrCodeInvalidFormat, Field: "date", Message: "expected YYYY-MM-DD format"}
	ErrInvalidLatitude    = &DomainError{Code: ErrCodeValidation, Field: "latitude", Message: "must be between -90 and 90"}
	ErrInvalidLongitude   = &DomainError{Code: ErrCodeValidation, Field: "longitude", Message: "must be between -180 and 180"}
	ErrAddressTooLong     = &DomainError{Code: ErrCodeLimitExceeded, Field: "address", Message: "address exceeds maximum length"}
	ErrDuplicateTag       = &DomainError{Code: ErrCodeDuplicate, Field: "tag", Message: "tag already exists"}
	ErrTagTooLong         = &DomainError{Code: ErrCodeLimitExceeded, Field: "tag", Message: "tag exceeds maximum length"}
	ErrEmptyTag           = &DomainError{Code: ErrCodeValidation, Field: "tag", Message: "tag cannot be empty"}
	ErrInvalidTagFormat   = &DomainError{Code: ErrCodeInvalidFormat, Field: "tag", Message: "tag contains invalid characters"}
	ErrMaxTagsExceeded    = &DomainError{Code: ErrCodeLimitExceeded, Field: "tags", Message: "maximum number of tags exceeded"}
	ErrModelNameTooLong   = &DomainError{Code: ErrCodeLimitExceeded, Field: "modelName", Message: "model name exceeds maximum length"}
	ErrEmptyModelName     = &DomainError{Code: ErrCodeValidation, Field: "modelName", Message: "model name cannot be empty when voice is enabled"}
	ErrNilNippou          = &DomainError{Code: ErrCodeNilReceiver, Field: "nippou", Message: "operation on nil Nippou"}
)

// IsValidationError checks if the error is a validation error.
func IsValidationError(err error) bool {
	var domErr *DomainError
	if errors.As(err, &domErr) {
		return domErr.Code == ErrCodeValidation
	}
	return false
}

// ============================================================================
// ID Value Object - Immutable Unique Identifier
// ============================================================================

// ID is an immutable Value Object representing a unique identifier for Nippou.
type ID struct {
	value string
}

// IDGenerator defines the interface for generating unique IDs (DIP compliance).
type IDGenerator interface {
	Generate() ID
}

// UUIDGenerator is the default implementation using UUID v4.
type UUIDGenerator struct{}

// Generate creates a new unique ID using UUID v4.
func (g *UUIDGenerator) Generate() ID {
	return ID{value: uuid.New().String()}
}

// defaultIDGenerator is the package-level default generator.
var defaultIDGenerator IDGenerator = &UUIDGenerator{}

// NewID generates a new unique ID using the default generator.
func NewID() ID {
	return defaultIDGenerator.Generate()
}

// IDFromString creates an ID from an existing string value.
// Returns error if the string is empty or invalid UUID format.
func IDFromString(s string) (ID, error) {
	if s == "" {
		return ID{}, &DomainError{Code: ErrCodeValidation, Field: "id", Message: "ID cannot be empty"}
	}
	// Validate UUID format
	if _, err := uuid.Parse(s); err != nil {
		return ID{}, &DomainError{Code: ErrCodeInvalidFormat, Field: "id", Message: "invalid UUID format"}
	}
	return ID{value: s}, nil
}

// String returns the string representation of the ID.
func (id ID) String() string {
	return id.value
}

// IsEmpty checks if the ID is empty.
func (id ID) IsEmpty() bool {
	return id.value == ""
}

// Equals compares two IDs for equality.
func (id ID) Equals(other ID) bool {
	return id.value == other.value
}

// ============================================================================
// Location Value Object - Immutable Geographical Coordinates
// ============================================================================

// Location is an immutable Value Object representing geographical coordinates.
type Location struct {
	latitude  float64
	longitude float64
	address   string
}

// NewLocation creates a new validated Location.
// Returns error if coordinates are out of valid range or address is too long.
func NewLocation(lat, lng float64, address string) (*Location, error) {
	if lat < MinLatitude || lat > MaxLatitude {
		return nil, ErrInvalidLatitude
	}
	if lng < MinLongitude || lng > MaxLongitude {
		return nil, ErrInvalidLongitude
	}
	sanitizedAddress := sanitizeString(address)
	if utf8.RuneCountInString(sanitizedAddress) > MaxAddressLength {
		return nil, ErrAddressTooLong
	}
	return &Location{
		latitude:  lat,
		longitude: lng,
		address:   sanitizedAddress,
	}, nil
}

// Latitude returns the latitude coordinate.
func (l *Location) Latitude() float64 {
	if l == nil {
		return 0
	}
	return l.latitude
}

// Longitude returns the longitude coordinate.
func (l *Location) Longitude() float64 {
	if l == nil {
		return 0
	}
	return l.longitude
}

// Address returns the sanitized address string.
func (l *Location) Address() string {
	if l == nil {
		return ""
	}
	return l.address
}

// Equals compares two Locations for equality.
func (l *Location) Equals(other *Location) bool {
	if l == nil && other == nil {
		return true
	}
	if l == nil || other == nil {
		return false
	}
	return l.latitude == other.latitude &&
		l.longitude == other.longitude &&
		l.address == other.address
}

// ============================================================================
// VoiceConfig Value Object - Immutable Voice Settings
// ============================================================================

// VoiceConfig is an immutable Value Object for voice settings.
type VoiceConfig struct {
	enabled   bool
	modelName string
}

// NewVoiceConfig creates a new validated VoiceConfig.
// Returns error if enabled is true but modelName is empty or too long.
func NewVoiceConfig(enabled bool, modelName string) (*VoiceConfig, error) {
	sanitizedModel := sanitizeString(modelName)
	if enabled {
		if sanitizedModel == "" {
			return nil, ErrEmptyModelName
		}
	}
	if utf8.RuneCountInString(sanitizedModel) > MaxModelNameLength {
		return nil, ErrModelNameTooLong
	}
	return &VoiceConfig{
		enabled:   enabled,
		modelName: sanitizedModel,
	}, nil
}

// Enabled returns whether voice is enabled.
func (v *VoiceConfig) Enabled() bool {
	if v == nil {
		return false
	}
	return v.enabled
}

// ModelName returns the voice model name.
func (v *VoiceConfig) ModelName() string {
	if v == nil {
		return ""
	}
	return v.modelName
}

// Equals compares two VoiceConfigs for equality.
func (v *VoiceConfig) Equals(other *VoiceConfig) bool {
	if v == nil && other == nil {
		return true
	}
	if v == nil || other == nil {
		return false
	}
	return v.enabled == other.enabled && v.modelName == other.modelName
}

// ============================================================================
// Tag Value Object - Validated Tag String
// ============================================================================

// tagPattern defines valid tag format: alphanumeric, hyphens, underscores.
var tagPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

// Tag is a validated, immutable tag value.
type Tag struct {
	value string
}

// NewTag creates a new validated Tag.
// Returns error if tag is empty, too long, or contains invalid characters.
func NewTag(value string) (Tag, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return Tag{}, ErrEmptyTag
	}
	if utf8.RuneCountInString(trimmed) > MaxTagLength {
		return Tag{}, ErrTagTooLong
	}
	// Convert to lowercase for case-insensitive comparison
	normalized := strings.ToLower(trimmed)
	if !tagPattern.MatchString(normalized) {
		return Tag{}, ErrInvalidTagFormat
	}
	return Tag{value: normalized}, nil
}

// String returns the tag value.
func (t Tag) String() string {
	return t.value
}

// Equals compares two tags (case-insensitive).
func (t Tag) Equals(other Tag) bool {
	return t.value == other.value
}

// ============================================================================
// Nippou Entity - Core Domain Entity
// ============================================================================

// Nippou is the core entity representing a daily report.
// All fields are private to ensure invariants are maintained.
type Nippou struct {
	id        ID
	date      time.Time
	content   string
	location  *Location
	voice     *VoiceConfig
	tags      []Tag
	createdAt time.Time
	updatedAt time.Time
}

// NippouBuilder provides a fluent API for creating Nippou entities.
type NippouBuilder struct {
	dateStr   string
	content   string
	location  *Location
	voice     *VoiceConfig
	tags      []Tag
	idGen     IDGenerator
	timeFunc  func() time.Time
}

// NewNippouBuilder creates a new builder with required fields.
func NewNippouBuilder(dateStr, content string) *NippouBuilder {
	return &NippouBuilder{
		dateStr:  dateStr,
		content:  content,
		tags:     make([]Tag, 0),
		idGen:    defaultIDGenerator,
		timeFunc: time.Now,
	}
}

// WithLocation sets the location.
func (b *NippouBuilder) WithLocation(location *Location) *NippouBuilder {
	b.location = location
	return b
}

// WithVoice sets the voice configuration.
func (b *NippouBuilder) WithVoice(voice *VoiceConfig) *NippouBuilder {
	b.voice = voice
	return b
}

// WithTags sets the initial tags.
func (b *NippouBuilder) WithTags(tags []Tag) *NippouBuilder {
	b.tags = tags
	return b
}

// WithIDGenerator sets a custom ID generator (useful for testing).
func (b *NippouBuilder) WithIDGenerator(gen IDGenerator) *NippouBuilder {
	b.idGen = gen
	return b
}

// WithTimeFunc sets a custom time function (useful for testing).
func (b *NippouBuilder) WithTimeFunc(fn func() time.Time) *NippouBuilder {
	b.timeFunc = fn
	return b
}

// Build creates the Nippou entity after validating all inputs.
func (b *NippouBuilder) Build() (*Nippou, error) {
	// Validate and sanitize content
	sanitizedContent := sanitizeString(b.content)
	if sanitizedContent == "" {
		return nil, ErrEmptyContent
	}
	if utf8.RuneCountInString(sanitizedContent) > MaxContentLength {
		return nil, ErrContentTooLong
	}

	// Parse and validate date
	date, err := time.Parse("2006-01-02", b.dateStr)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	// Validate tag count
	if len(b.tags) > MaxTagCount {
		return nil, ErrMaxTagsExceeded
	}

	now := b.timeFunc()
	return &Nippou{
		id:        b.idGen.Generate(),
		date:      date,
		content:   sanitizedContent,
		location:  b.location,
		voice:     b.voice,
		tags:      copyTags(b.tags),
		createdAt: now,
		updatedAt: now,
	}, nil
}

// NewNippou creates a new Nippou entity with the given date and content.
// Date must be in "YYYY-MM-DD" format. This is a convenience function
// that uses the builder internally.
func NewNippou(dateStr string, content string) (*Nippou, error) {
	return NewNippouBuilder(dateStr, content).Build()
}

// ============================================================================
// Nippou Getters - Safe Read Access
// ============================================================================

// ID returns the unique identifier.
func (n *Nippou) ID() ID {
	if n == nil {
		return ID{}
	}
	return n.id
}

// Date returns the report date.
func (n *Nippou) Date() time.Time {
	if n == nil {
		return time.Time{}
	}
	return n.date
}

// Content returns the sanitized content.
func (n *Nippou) Content() string {
	if n == nil {
		return ""
	}
	return n.content
}

// Location returns a copy of the location (nil-safe).
func (n *Nippou) Location() *Location {
	if n == nil || n.location == nil {
		return nil
	}
	// Return a copy to prevent external mutation
	return &Location{
		latitude:  n.location.latitude,
		longitude: n.location.longitude,
		address:   n.location.address,
	}
}

// Voice returns a copy of the voice config (nil-safe).
func (n *Nippou) Voice() *VoiceConfig {
	if n == nil || n.voice == nil {
		return nil
	}
	// Return a copy to prevent external mutation
	return &VoiceConfig{
		enabled:   n.voice.enabled,
		modelName: n.voice.modelName,
	}
}

// Tags returns a copy of all tags.
func (n *Nippou) Tags() []Tag {
	if n == nil {
		return nil
	}
	return copyTags(n.tags)
}

// TagStrings returns all tags as strings.
func (n *Nippou) TagStrings() []string {
	if n == nil {
		return nil
	}
	result := make([]string, len(n.tags))
	for i, t := range n.tags {
		result[i] = t.String()
	}
	return result
}

// CreatedAt returns the creation timestamp.
func (n *Nippou) CreatedAt() time.Time {
	if n == nil {
		return time.Time{}
	}
	return n.createdAt
}

// UpdatedAt returns the last update timestamp.
func (n *Nippou) UpdatedAt() time.Time {
	if n == nil {
		return time.Time{}
	}
	return n.updatedAt
}

// ============================================================================
// Nippou Mutators - Safe Write Operations
// ============================================================================

// UpdateContent updates the content with validation.
func (n *Nippou) UpdateContent(content string) error {
	if n == nil {
		return ErrNilNippou
	}
	sanitized := sanitizeString(content)
	if sanitized == "" {
		return ErrEmptyContent
	}
	if utf8.RuneCountInString(sanitized) > MaxContentLength {
		return ErrContentTooLong
	}
	n.content = sanitized
	n.updatedAt = time.Now()
	return nil
}

// AttachLocation attaches a validated location to the Nippou.
func (n *Nippou) AttachLocation(lat, lng float64, address string) error {
	if n == nil {
		return ErrNilNippou
	}
	loc, err := NewLocation(lat, lng, address)
	if err != nil {
		return err
	}
	n.location = loc
	n.updatedAt = time.Now()
	return nil
}

// SetLocation sets a pre-validated location.
func (n *Nippou) SetLocation(loc *Location) error {
	if n == nil {
		return ErrNilNippou
	}
	n.location = loc
	n.updatedAt = time.Now()
	return nil
}

// RemoveLocation removes the attached location.
func (n *Nippou) RemoveLocation() error {
	if n == nil {
		return ErrNilNippou
	}
	n.location = nil
	n.updatedAt = time.Now()
	return nil
}

// SetVoiceConfig sets the voice configuration with validation.
func (n *Nippou) SetVoiceConfig(enabled bool, modelName string) error {
	if n == nil {
		return ErrNilNippou
	}
	voice, err := NewVoiceConfig(enabled, modelName)
	if err != nil {
		return err
	}
	n.voice = voice
	n.updatedAt = time.Now()
	return nil
}

// SetVoice sets a pre-validated voice config.
func (n *Nippou) SetVoice(voice *VoiceConfig) error {
	if n == nil {
		return ErrNilNippou
	}
	n.voice = voice
	n.updatedAt = time.Now()
	return nil
}

// RemoveVoice removes the voice configuration.
func (n *Nippou) RemoveVoice() error {
	if n == nil {
		return ErrNilNippou
	}
	n.voice = nil
	n.updatedAt = time.Now()
	return nil
}

// AddTag adds a validated tag to the Nippou.
func (n *Nippou) AddTag(tagStr string) error {
	if n == nil {
		return ErrNilNippou
	}
	tag, err := NewTag(tagStr)
	if err != nil {
		return err
	}
	return n.AddValidatedTag(tag)
}

// AddValidatedTag adds a pre-validated tag.
func (n *Nippou) AddValidatedTag(tag Tag) error {
	if n == nil {
		return ErrNilNippou
	}
	if len(n.tags) >= MaxTagCount {
		return ErrMaxTagsExceeded
	}
	for _, t := range n.tags {
		if t.Equals(tag) {
			return ErrDuplicateTag
		}
	}
	n.tags = append(n.tags, tag)
	n.updatedAt = time.Now()
	return nil
}

// RemoveTag removes a tag by value.
func (n *Nippou) RemoveTag(tagStr string) error {
	if n == nil {
		return ErrNilNippou
	}
	tag, err := NewTag(tagStr)
	if err != nil {
		return err
	}
	for i, t := range n.tags {
		if t.Equals(tag) {
			n.tags = append(n.tags[:i], n.tags[i+1:]...)
			n.updatedAt = time.Now()
			return nil
		}
	}
	return nil // Not found is not an error
}

// HasTag checks if the Nippou has the specified tag (case-insensitive).
func (n *Nippou) HasTag(tagStr string) bool {
	if n == nil {
		return false
	}
	tag, err := NewTag(tagStr)
	if err != nil {
		return false
	}
	for _, t := range n.tags {
		if t.Equals(tag) {
			return true
		}
	}
	return false
}

// TagCount returns the number of tags.
func (n *Nippou) TagCount() int {
	if n == nil {
		return 0
	}
	return len(n.tags)
}

// ============================================================================
// Repository Interface - Persistence Abstraction (ISP & DIP)
// ============================================================================

// Reader defines read operations for Nippou persistence.
type Reader interface {
	FindByID(id ID) (*Nippou, error)
	FindByDate(date time.Time) ([]*Nippou, error)
}

// Writer defines write operations for Nippou persistence.
type Writer interface {
	Save(n *Nippou) error
	Delete(id ID) error
}

// Repository combines Reader and Writer interfaces.
// Implementation is provided in the Infrastructure layer.
type Repository interface {
	Reader
	Writer
}

// ============================================================================
// Reconstruction - For Repository Implementation
// ============================================================================

// ReconstructedNippou contains all fields needed to reconstruct a Nippou from storage.
type ReconstructedNippou struct {
	ID        string
	Date      time.Time
	Content   string
	Location  *Location
	Voice     *VoiceConfig
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Reconstruct creates a Nippou from stored data without validation.
// This should only be used by repository implementations.
func Reconstruct(data ReconstructedNippou) (*Nippou, error) {
	id, err := IDFromString(data.ID)
	if err != nil {
		return nil, err
	}

	tags := make([]Tag, 0, len(data.Tags))
	for _, tagStr := range data.Tags {
		tag, err := NewTag(tagStr)
		if err != nil {
			// Skip invalid tags during reconstruction
			continue
		}
		tags = append(tags, tag)
	}

	return &Nippou{
		id:        id,
		date:      data.Date,
		content:   data.Content,
		location:  data.Location,
		voice:     data.Voice,
		tags:      tags,
		createdAt: data.CreatedAt,
		updatedAt: data.UpdatedAt,
	}, nil
}

// ============================================================================
// Helper Functions - Internal Utilities
// ============================================================================

// sanitizeString removes potentially dangerous characters and trims whitespace.
func sanitizeString(s string) string {
	// Trim whitespace
	s = strings.TrimSpace(s)
	// Remove null bytes (potential security issue)
	s = strings.ReplaceAll(s, "\x00", "")
	// Normalize line endings to LF
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

// copyTags creates a defensive copy of a tag slice.
func copyTags(tags []Tag) []Tag {
	if tags == nil {
		return nil
	}
	result := make([]Tag, len(tags))
	copy(result, tags)
	return result
}
