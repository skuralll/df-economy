package service

import (
	"errors"
	"fmt"
)

var (
	ErrPlayerExists = errors.New("player already exists")
	ErrValidation   = errors.New("validation error")
)

func NewPlayerExistsError(id string) error {
	return fmt.Errorf("%w: player ID %s already exists", ErrPlayerExists, id)
}

func NewValidationError(field, message string) error {
	return fmt.Errorf("%w: %s %s", ErrValidation, field, message)
}
