package service

import (
	"errors"
	"fmt"
)

var (
	ErrPlayerExists   = errors.New("player already exists")
	ErrValidation     = errors.New("validation error")
	ErrUnknownPlayer  = errors.New("unknown player")
	ErrInternalError  = errors.New("internal error")
)

func NewPlayerExistsError(id string) error {
	return fmt.Errorf("%w: player ID %s already exists", ErrPlayerExists, id)
}

func NewValidationError(field, message string) error {
	return fmt.Errorf("%w: %s %s", ErrValidation, field, message)
}

func NewUnknownPlayerError(identifier string) error {
	return fmt.Errorf("%w: %s", ErrUnknownPlayer, identifier)
}

func NewInternalError(operation, message string) error {
	return fmt.Errorf("%w: %s failed: %s", ErrInternalError, operation, message)
}
