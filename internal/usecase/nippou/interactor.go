package nippou

import (
	"context"

	domain "salesforce-mcp-server/internal/domain/nippou"
)

// ============================================================================
// Create UseCase - Application Service
// ============================================================================

// CreateUseCase handles the creation of Nippou entities.
type CreateUseCase struct {
	repo domain.Repository
}

// NewCreateUseCase creates a new CreateUseCase with the given repository.
func NewCreateUseCase(repo domain.Repository) (*CreateUseCase, error) {
	if repo == nil {
		return nil, ErrRepositoryNil
	}
	return &CreateUseCase{repo: repo}, nil
}

// Execute creates a new Nippou based on the input.
// It validates input, creates the domain entity, persists it, and returns the output DTO.
func (uc *CreateUseCase) Execute(ctx context.Context, input *CreateInput) (*CreateOutput, error) {
	// Guard: Check context
	if ctx == nil {
		return nil, ErrContextNil
	}

	// Check for context cancellation early
	select {
	case <-ctx.Done():
		return nil, &UseCaseError{
			Code:    ErrCodeContextCancelled,
			Message: "operation cancelled",
			Cause:   ctx.Err(),
		}
	default:
	}

	// Step 1: Validate input DTO
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// Step 2: Build domain entity using the builder
	builder := domain.NewNippouBuilder(input.Date, input.Content)

	// Step 3: Add optional location
	if input.Location != nil {
		loc, err := domain.NewLocation(
			input.Location.Latitude,
			input.Location.Longitude,
			input.Location.Address,
		)
		if err != nil {
			return nil, NewDomainViolationError(err)
		}
		builder.WithLocation(loc)
	}

	// Step 4: Add optional voice config
	if input.Voice != nil {
		voice, err := domain.NewVoiceConfig(input.Voice.Enabled, input.Voice.ModelName)
		if err != nil {
			return nil, NewDomainViolationError(err)
		}
		builder.WithVoice(voice)
	}

	// Step 5: Add tags
	if len(input.Tags) > 0 {
		tags := make([]domain.Tag, 0, len(input.Tags))
		for _, tagStr := range input.Tags {
			tag, err := domain.NewTag(tagStr)
			if err != nil {
				return nil, NewDomainViolationError(err)
			}
			// Check for duplicate tags
			for _, existingTag := range tags {
				if existingTag.Equals(tag) {
					return nil, NewDomainViolationError(domain.ErrDuplicateTag)
				}
			}
			tags = append(tags, tag)
		}
		builder.WithTags(tags)
	}

	// Step 6: Build the entity
	nippou, err := builder.Build()
	if err != nil {
		return nil, NewDomainViolationError(err)
	}

	// Check context before repository operation
	select {
	case <-ctx.Done():
		return nil, &UseCaseError{
			Code:    ErrCodeContextCancelled,
			Message: "operation cancelled before persistence",
			Cause:   ctx.Err(),
		}
	default:
	}

	// Step 7: Persist to repository
	if err := uc.repo.Save(nippou); err != nil {
		return nil, NewRepositoryError(err)
	}

	// Step 8: Map to output DTO
	return mapToOutput(nippou), nil
}

// mapToOutput converts a domain Nippou to the output DTO.
func mapToOutput(n *domain.Nippou) *CreateOutput {
	if n == nil {
		return nil
	}

	output := &CreateOutput{
		ID:        n.ID().String(),
		Date:      n.Date().Format("2006-01-02"),
		Content:   n.Content(),
		Tags:      n.TagStrings(),
		CreatedAt: n.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: n.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	// Ensure tags is never nil in output
	if output.Tags == nil {
		output.Tags = []string{}
	}

	// Map location if present
	if loc := n.Location(); loc != nil {
		output.Location = &LocationOutput{
			Latitude:  loc.Latitude(),
			Longitude: loc.Longitude(),
			Address:   loc.Address(),
		}
	}

	// Map voice if present
	if voice := n.Voice(); voice != nil {
		output.Voice = &VoiceOutput{
			Enabled:   voice.Enabled(),
			ModelName: voice.ModelName(),
		}
	}

	return output
}
