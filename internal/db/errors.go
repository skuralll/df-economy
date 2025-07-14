package db

import (
	"errors"
	"fmt"
)

// Sentinel errors for DB layer
var (
	ErrValidation          = errors.New("validation error")
	ErrDatabase            = errors.New("database error")
	ErrNotFound            = errors.New("record not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

// NewValidationError creates a new validation error with detailed information
func NewValidationError(field, message string) error {
	return fmt.Errorf("%w: %s %s", ErrValidation, field, message)
}

// NewDatabaseError creates a new database error with operation context
func NewDatabaseError(operation, message string) error {
	return fmt.Errorf("%w: %s failed: %s", ErrDatabase, operation, message)
}

// NewNotFoundError creates a new not found error with resource context
func NewNotFoundError(resource string) error {
	return fmt.Errorf("%w: %s not found", ErrNotFound, resource)
}

// NewInsufficientBalanceError creates a new insufficient balance error with amount details
func NewInsufficientBalanceError(required, available float64) error {
	return fmt.Errorf("%w: required %.2f, available %.2f", ErrInsufficientBalance, required, available)
}