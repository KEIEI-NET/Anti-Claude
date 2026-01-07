package nippou

import (
	"fmt"
	"strings"
	"unicode/utf8"

	domain "salesforce-mcp-server/internal/domain/nippou"
)

// ============================================================================
// Input DTO - Request Data Transfer Object
// ============================================================================

// LocationInput represents location data in the request.
type LocationInput struct {
	Latitude  float64
	Longitude float64
	Address   string
}

// VoiceInput represents voice configuration in the request.
type VoiceInput struct {
	Enabled   bool
	ModelName string
}

// CreateInput is the input DTO for creating a Nippou.
type CreateInput struct {
	Date     string
	Content  string
	Location *LocationInput
	Voice    *VoiceInput
	Tags     []string
}

// Validate performs early validation on the input DTO.
// This catches obvious errors before domain processing.
func (i *CreateInput) Validate() error {
	if i == nil {
		return ErrNilInput
	}

	// Validate required fields
	if strings.TrimSpace(i.Date) == "" {
		return NewInvalidInputError("date", "cannot be empty")
	}
	if strings.TrimSpace(i.Content) == "" {
		return NewInvalidInputError("content", "cannot be empty")
	}

	// Validate content length (early check to avoid processing large payloads)
	if utf8.RuneCountInString(i.Content) > domain.MaxContentLength {
		return NewInvalidInputError("content", fmt.Sprintf("exceeds maximum length of %d characters", domain.MaxContentLength))
	}

	// Validate tags count
	if len(i.Tags) > domain.MaxTagCount {
		return NewInvalidInputError("tags", fmt.Sprintf("exceeds maximum count of %d", domain.MaxTagCount))
	}

	// Validate individual tag lengths
	for idx, tag := range i.Tags {
		if utf8.RuneCountInString(tag) > domain.MaxTagLength {
			return NewInvalidInputError("tags", fmt.Sprintf("tag at index %d exceeds maximum length of %d", idx, domain.MaxTagLength))
		}
	}

	// Validate location if provided
	if i.Location != nil {
		if i.Location.Latitude < domain.MinLatitude || i.Location.Latitude > domain.MaxLatitude {
			return NewInvalidInputError("location.latitude", "must be between -90 and 90")
		}
		if i.Location.Longitude < domain.MinLongitude || i.Location.Longitude > domain.MaxLongitude {
			return NewInvalidInputError("location.longitude", "must be between -180 and 180")
		}
		if utf8.RuneCountInString(i.Location.Address) > domain.MaxAddressLength {
			return NewInvalidInputError("location.address", fmt.Sprintf("exceeds maximum length of %d", domain.MaxAddressLength))
		}
	}

	// Validate voice if provided
	if i.Voice != nil {
		if i.Voice.Enabled && strings.TrimSpace(i.Voice.ModelName) == "" {
			return NewInvalidInputError("voice.modelName", "cannot be empty when voice is enabled")
		}
		if utf8.RuneCountInString(i.Voice.ModelName) > domain.MaxModelNameLength {
			return NewInvalidInputError("voice.modelName", fmt.Sprintf("exceeds maximum length of %d", domain.MaxModelNameLength))
		}
	}

	return nil
}

// ============================================================================
// Output DTO - Response Data Transfer Object
// ============================================================================

// LocationOutput represents location data in the response.
type LocationOutput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

// VoiceOutput represents voice configuration in the response.
type VoiceOutput struct {
	Enabled   bool   `json:"enabled"`
	ModelName string `json:"modelName"`
}

// CreateOutput is the output DTO for the created Nippou.
type CreateOutput struct {
	ID        string          `json:"id"`
	Date      string          `json:"date"`
	Content   string          `json:"content"`
	Location  *LocationOutput `json:"location,omitempty"`
	Voice     *VoiceOutput    `json:"voice,omitempty"`
	Tags      []string        `json:"tags"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}
