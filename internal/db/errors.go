package db

import (
	"errors"
	"fmt"
)

// Sentinel errors for DB layer
var (
	ErrValidation = errors.New("validation error")
	ErrDatabase   = errors.New("database error")
)

// NewValidationError creates a new validation error with detailed information
func NewValidationError(field, message string) error {
	return fmt.Errorf("%w: %s %s", ErrValidation, field, message)
}

// NewDatabaseError creates a new database error with operation context
func NewDatabaseError(operation, message string) error {
	return fmt.Errorf("%w: %s failed: %s", ErrDatabase, operation, message)
}