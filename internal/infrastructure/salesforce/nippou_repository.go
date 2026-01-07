package salesforce

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"salesforce-mcp-server/internal/domain/nippou"
)

// ============================================================================
// Repository Errors - Infrastructure Layer Error Handling
// ============================================================================

// RepositoryError represents a repository-level error with context.
type RepositoryError struct {
	Operation string
	Cause     error
}

func (e *RepositoryError) Error() string {
	return fmt.Sprintf("repository %s failed: %v", e.Operation, e.Cause)
}

func (e *RepositoryError) Unwrap() error {
	return e.Cause
}

// ============================================================================
// NippouRepository - Implements domain.Repository
// ============================================================================

// NippouRepository implements nippou.Repository interface using Salesforce as backend.
type NippouRepository struct {
	client *Client
	ctx    context.Context
}

// NewNippouRepository creates a new NippouRepository with the given Salesforce client.
func NewNippouRepository(client *Client) *NippouRepository {
	return &NippouRepository{
		client: client,
		ctx:    context.Background(),
	}
}

// NewNippouRepositoryWithContext creates a repository with a custom context.
func NewNippouRepositoryWithContext(ctx context.Context, client *Client) *NippouRepository {
	return &NippouRepository{
		client: client,
		ctx:    ctx,
	}
}

// WithContext returns a new repository with the given context.
func (r *NippouRepository) WithContext(ctx context.Context) *NippouRepository {
	return &NippouRepository{
		client: r.client,
		ctx:    ctx,
	}
}

// ============================================================================
// Reader Interface Implementation
// ============================================================================

// FindByID retrieves a Nippou by its unique ID.
// Returns nil if the record is not found.
func (r *NippouRepository) FindByID(id nippou.ID) (*nippou.Nippou, error) {
	if id.IsEmpty() {
		return nil, &RepositoryError{
			Operation: "FindByID",
			Cause:     fmt.Errorf("empty ID provided"),
		}
	}

	// Query by ID using SOQL for consistency
	// (Salesforce ID != our UUID, so we store our ID in a custom field or use it as external ID)
	soql := fmt.Sprintf(
		"SELECT Id, Date__c, Content__c, Latitude__c, Longitude__c, Address__c, "+
			"VoiceEnabled__c, VoiceModel__c, Tags__c, CreatedDate, LastModifiedDate "+
			"FROM %s WHERE Id = '%s' LIMIT 1",
		NippouObjectName,
		EscapeSOQL(id.String()),
	)

	var result QueryResult
	if err := r.client.Query(r.ctx, url.QueryEscape(soql), &result); err != nil {
		// Check if it's a not found error
		if apiErr, ok := err.(*APIError); ok && apiErr.IsNotFound() {
			return nil, nil
		}
		return nil, &RepositoryError{
			Operation: "FindByID",
			Cause:     err,
		}
	}

	if result.TotalSize == 0 || len(result.Records) == 0 {
		return nil, nil
	}

	n, err := result.Records[0].ToDomain()
	if err != nil {
		return nil, &RepositoryError{
			Operation: "FindByID",
			Cause:     fmt.Errorf("failed to convert to domain: %w", err),
		}
	}

	return n, nil
}

// FindByDate retrieves all Nippou entries for a specific date.
func (r *NippouRepository) FindByDate(date time.Time) ([]*nippou.Nippou, error) {
	dateStr := FormatDateForSOQL(date)

	soql := fmt.Sprintf(
		"SELECT Id, Date__c, Content__c, Latitude__c, Longitude__c, Address__c, "+
			"VoiceEnabled__c, VoiceModel__c, Tags__c, CreatedDate, LastModifiedDate "+
			"FROM %s WHERE Date__c = %s ORDER BY CreatedDate ASC",
		NippouObjectName,
		dateStr,
	)

	return r.executeQuery(soql, "FindByDate")
}

// ============================================================================
// Writer Interface Implementation
// ============================================================================

// Save persists a Nippou entity to Salesforce.
// If the entity has no Salesforce ID, it creates a new record.
// If it has an ID, it updates the existing record.
func (r *NippouRepository) Save(n *nippou.Nippou) error {
	if n == nil {
		return &RepositoryError{
			Operation: "Save",
			Cause:     fmt.Errorf("nil Nippou provided"),
		}
	}

	sfRecord := FromDomain(n)

	// Check if this is a new record or an update
	if n.ID().IsEmpty() {
		return &RepositoryError{
			Operation: "Save",
			Cause:     fmt.Errorf("Nippou must have a valid ID"),
		}
	}

	// Try to find existing record first
	existing, err := r.FindByID(n.ID())
	if err != nil {
		return &RepositoryError{
			Operation: "Save",
			Cause:     fmt.Errorf("failed to check existing record: %w", err),
		}
	}

	if existing == nil {
		// Create new record
		return r.create(sfRecord)
	}

	// Update existing record
	return r.update(sfRecord.ID, sfRecord)
}

// Delete removes a Nippou by its ID.
func (r *NippouRepository) Delete(id nippou.ID) error {
	if id.IsEmpty() {
		return &RepositoryError{
			Operation: "Delete",
			Cause:     fmt.Errorf("empty ID provided"),
		}
	}

	err := r.client.DeleteSObject(r.ctx, NippouObjectName, id.String())
	if err != nil {
		// Not found is not an error for delete operations
		if apiErr, ok := err.(*APIError); ok && apiErr.IsNotFound() {
			return nil
		}
		return &RepositoryError{
			Operation: "Delete",
			Cause:     err,
		}
	}

	return nil
}

// ============================================================================
// Additional Query Methods (Extended Interface)
// ============================================================================

// FindByDateRange retrieves all Nippou entries within a date range.
func (r *NippouRepository) FindByDateRange(startDate, endDate time.Time) ([]*nippou.Nippou, error) {
	startStr := FormatDateForSOQL(startDate)
	endStr := FormatDateForSOQL(endDate)

	soql := fmt.Sprintf(
		"SELECT Id, Date__c, Content__c, Latitude__c, Longitude__c, Address__c, "+
			"VoiceEnabled__c, VoiceModel__c, Tags__c, CreatedDate, LastModifiedDate "+
			"FROM %s WHERE Date__c >= %s AND Date__c <= %s ORDER BY Date__c ASC, CreatedDate ASC",
		NippouObjectName,
		startStr,
		endStr,
	)

	return r.executeQuery(soql, "FindByDateRange")
}

// FindByTag retrieves all Nippou entries that contain a specific tag.
func (r *NippouRepository) FindByTag(tag string) ([]*nippou.Nippou, error) {
	// Using LIKE for substring match in comma-separated tags
	escapedTag := EscapeSOQL(tag)

	soql := fmt.Sprintf(
		"SELECT Id, Date__c, Content__c, Latitude__c, Longitude__c, Address__c, "+
			"VoiceEnabled__c, VoiceModel__c, Tags__c, CreatedDate, LastModifiedDate "+
			"FROM %s WHERE Tags__c LIKE '%%%s%%' ORDER BY CreatedDate DESC",
		NippouObjectName,
		escapedTag,
	)

	return r.executeQuery(soql, "FindByTag")
}

// ============================================================================
// Internal Helper Methods
// ============================================================================

// create inserts a new record in Salesforce.
func (r *NippouRepository) create(sf *NippouSF) error {
	payload := sf.ToCreatePayload()

	result, err := r.client.CreateSObject(r.ctx, NippouObjectName, payload)
	if err != nil {
		return &RepositoryError{
			Operation: "Save(create)",
			Cause:     err,
		}
	}

	if !result.Success {
		return &RepositoryError{
			Operation: "Save(create)",
			Cause:     fmt.Errorf("Salesforce create returned success=false"),
		}
	}

	return nil
}

// update modifies an existing record in Salesforce.
func (r *NippouRepository) update(sfID string, sf *NippouSF) error {
	payload := sf.ToUpdatePayload()

	err := r.client.UpdateSObject(r.ctx, NippouObjectName, sfID, payload)
	if err != nil {
		return &RepositoryError{
			Operation: "Save(update)",
			Cause:     err,
		}
	}

	return nil
}

// executeQuery runs a SOQL query and converts results to domain entities.
func (r *NippouRepository) executeQuery(soql, operation string) ([]*nippou.Nippou, error) {
	var result QueryResult
	if err := r.client.Query(r.ctx, url.QueryEscape(soql), &result); err != nil {
		return nil, &RepositoryError{
			Operation: operation,
			Cause:     err,
		}
	}

	nippous := make([]*nippou.Nippou, 0, len(result.Records))
	for _, record := range result.Records {
		n, err := record.ToDomain()
		if err != nil {
			// Log warning but continue with other records
			continue
		}
		nippous = append(nippous, n)
	}

	return nippous, nil
}

// ============================================================================
// Compile-time Interface Compliance Check
// ============================================================================

// Ensure NippouRepository implements nippou.Repository at compile time.
var _ nippou.Repository = (*NippouRepository)(nil)
