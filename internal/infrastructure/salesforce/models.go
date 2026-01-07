package salesforce

import (
	"strings"
	"time"

	"salesforce-mcp-server/internal/domain/nippou"
)

// ============================================================================
// Salesforce Object Names
// ============================================================================

const (
	// NippouObjectName is the Salesforce custom object API name for Nippou.
	NippouObjectName = "Nippou__c"
)

// ============================================================================
// Nippou__c - Salesforce Custom Object Mapping
// ============================================================================

// NippouSF represents the Salesforce Nippou__c custom object structure.
// Field names follow Salesforce custom field naming convention (suffixed with __c).
type NippouSF struct {
	// Standard Salesforce fields
	ID string `json:"Id,omitempty"`

	// Custom fields for Nippou__c
	Date      string  `json:"Date__c,omitempty"`       // Date in YYYY-MM-DD format
	Content   string  `json:"Content__c,omitempty"`    // Long text area
	Latitude  float64 `json:"Latitude__c,omitempty"`   // Decimal (10,7)
	Longitude float64 `json:"Longitude__c,omitempty"`  // Decimal (10,7)
	Address   string  `json:"Address__c,omitempty"`    // Text(500)
	VoiceOn   bool    `json:"VoiceEnabled__c"`         // Checkbox
	VoiceModel string `json:"VoiceModel__c,omitempty"` // Text(100)
	Tags      string  `json:"Tags__c,omitempty"`       // Long text (comma-separated)

	// Audit fields (read-only from SF)
	CreatedDate      string `json:"CreatedDate,omitempty"`
	LastModifiedDate string `json:"LastModifiedDate,omitempty"`
}

// ============================================================================
// Domain <-> Salesforce Mapping Functions
// ============================================================================

// FromDomain converts a domain Nippou entity to Salesforce NippouSF structure.
// This is used when creating or updating records in Salesforce.
func FromDomain(n *nippou.Nippou) *NippouSF {
	if n == nil {
		return nil
	}

	sf := &NippouSF{
		ID:      n.ID().String(),
		Date:    n.Date().Format("2006-01-02"),
		Content: n.Content(),
	}

	// Map location if present
	if loc := n.Location(); loc != nil {
		sf.Latitude = loc.Latitude()
		sf.Longitude = loc.Longitude()
		sf.Address = loc.Address()
	}

	// Map voice config if present
	if voice := n.Voice(); voice != nil {
		sf.VoiceOn = voice.Enabled()
		sf.VoiceModel = voice.ModelName()
	}

	// Map tags as comma-separated string
	if tags := n.TagStrings(); len(tags) > 0 {
		sf.Tags = strings.Join(tags, ",")
	}

	return sf
}

// ToCreatePayload returns the NippouSF structure suitable for POST (create).
// Excludes ID and audit fields that Salesforce generates.
func (sf *NippouSF) ToCreatePayload() map[string]interface{} {
	payload := make(map[string]interface{})

	if sf.Date != "" {
		payload["Date__c"] = sf.Date
	}
	if sf.Content != "" {
		payload["Content__c"] = sf.Content
	}
	if sf.Latitude != 0 || sf.Longitude != 0 {
		payload["Latitude__c"] = sf.Latitude
		payload["Longitude__c"] = sf.Longitude
	}
	if sf.Address != "" {
		payload["Address__c"] = sf.Address
	}
	payload["VoiceEnabled__c"] = sf.VoiceOn
	if sf.VoiceModel != "" {
		payload["VoiceModel__c"] = sf.VoiceModel
	}
	if sf.Tags != "" {
		payload["Tags__c"] = sf.Tags
	}

	return payload
}

// ToUpdatePayload returns the NippouSF structure suitable for PATCH (update).
// Excludes ID, audit fields, and unchanged fields for efficiency.
func (sf *NippouSF) ToUpdatePayload() map[string]interface{} {
	// For simplicity, we update all fields (Salesforce handles null gracefully)
	return sf.ToCreatePayload()
}

// ToDomain converts a Salesforce NippouSF to domain Nippou entity.
// This is used when reading records from Salesforce.
func (sf *NippouSF) ToDomain() (*nippou.Nippou, error) {
	if sf == nil {
		return nil, nil
	}

	// Parse date
	date, err := time.Parse("2006-01-02", sf.Date)
	if err != nil {
		// Try ISO 8601 format that Salesforce sometimes returns
		date, err = time.Parse(time.RFC3339, sf.Date)
		if err != nil {
			return nil, err
		}
	}

	// Parse audit timestamps
	createdAt := parseTimestamp(sf.CreatedDate)
	updatedAt := parseTimestamp(sf.LastModifiedDate)
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}

	// Build location if coordinates are present
	var location *nippou.Location
	if sf.Latitude != 0 || sf.Longitude != 0 || sf.Address != "" {
		loc, err := nippou.NewLocation(sf.Latitude, sf.Longitude, sf.Address)
		if err == nil {
			location = loc
		}
	}

	// Build voice config
	var voice *nippou.VoiceConfig
	if sf.VoiceOn || sf.VoiceModel != "" {
		v, err := nippou.NewVoiceConfig(sf.VoiceOn, sf.VoiceModel)
		if err == nil {
			voice = v
		}
	}

	// Parse tags from comma-separated string
	var tags []string
	if sf.Tags != "" {
		for _, t := range strings.Split(sf.Tags, ",") {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	// Use Reconstruct to create domain entity from stored data
	return nippou.Reconstruct(nippou.ReconstructedNippou{
		ID:        sf.ID,
		Date:      date,
		Content:   sf.Content,
		Location:  location,
		Voice:     voice,
		Tags:      tags,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
}

// ============================================================================
// SOQL Query Result Structures
// ============================================================================

// QueryResult represents a SOQL query response from Salesforce.
type QueryResult struct {
	TotalSize      int         `json:"totalSize"`
	Done           bool        `json:"done"`
	NextRecordsURL string      `json:"nextRecordsUrl,omitempty"`
	Records        []NippouSF  `json:"records"`
}

// ============================================================================
// Helper Functions
// ============================================================================

// parseTimestamp parses various Salesforce timestamp formats.
func parseTimestamp(s string) time.Time {
	if s == "" {
		return time.Time{}
	}

	// Try ISO 8601 with timezone (Salesforce standard)
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}

	// Try ISO 8601 without timezone
	if t, err := time.Parse("2006-01-02T15:04:05.000+0000", s); err == nil {
		return t
	}

	// Try date only
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t
	}

	return time.Time{}
}

// EscapeSOQL escapes special characters in SOQL string literals.
// This prevents SOQL injection attacks.
func EscapeSOQL(s string) string {
	// Escape single quotes by doubling them
	s = strings.ReplaceAll(s, "'", "''")
	// Escape backslashes
	s = strings.ReplaceAll(s, "\\", "\\\\")
	return s
}

// FormatDateForSOQL formats a Go time.Time for SOQL queries.
func FormatDateForSOQL(t time.Time) string {
	return t.Format("2006-01-02")
}
