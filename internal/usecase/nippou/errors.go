package nippou

import (
	"errors"
	"fmt"
)

// ============================================================================
// UseCase Errors - Application Layer Error Handling
// ============================================================================

// UseCaseError represents an application-level error with code and context.
type UseCaseError struct {
	Code    string
	Message string
	Cause   error
}

func (e *UseCaseError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *UseCaseError) Unwrap() error {
	return e.Cause
}

// UseCase error codes.
const (
	ErrCodeInvalidInput     = "INVALID_INPUT"
	ErrCodeRepositoryError  = "REPOSITORY_ERROR"
	ErrCodeDomainViolation  = "DOMAIN_VIOLATION"
	ErrCodeContextCancelled = "CONTEXT_CANCELLED"
)

// Predefined usecase errors.
var (
	ErrNilInput      = &UseCaseError{Code: ErrCodeInvalidInput, Message: "input cannot be nil"}
	ErrContextNil    = &UseCaseError{Code: ErrCodeInvalidInput, Message: "context cannot be nil"}
	ErrRepositoryNil = &UseCaseError{Code: ErrCodeInvalidInput, Message: "repository cannot be nil"}
)

// NewInvalidInputError creates an input validation error.
func NewInvalidInputError(field, message string) *UseCaseError {
	return &UseCaseError{
		Code:    ErrCodeInvalidInput,
		Message: fmt.Sprintf("%s: %s", field, message),
	}
}

// NewRepositoryError wraps a repository error.
func NewRepositoryError(cause error) *UseCaseError {
	return &UseCaseError{
		Code:    ErrCodeRepositoryError,
		Message: "failed to persist nippou",
		Cause:   cause,
	}
}

// NewDomainViolationError wraps a domain error.
func NewDomainViolationError(cause error) *UseCaseError {
	return &UseCaseError{
		Code:    ErrCodeDomainViolation,
		Message: "domain validation failed",
		Cause:   cause,
	}
}

// IsUseCaseError checks if the error is a UseCaseError.
func IsUseCaseError(err error) bool {
	var ucErr *UseCaseError
	return errors.As(err, &ucErr)
}

// IsDomainViolation checks if the error is a domain violation.
func IsDomainViolation(err error) bool {
	var ucErr *UseCaseError
	if errors.As(err, &ucErr) {
		return ucErr.Code == ErrCodeDomainViolation
	}
	return false
}
