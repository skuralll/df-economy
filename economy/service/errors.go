package service

import (
	"errors"
	"fmt"
)

var (
	ErrPlayerExists = errors.New("player already exists")
)

func NewPlayerExistsError(id string) error {
	return fmt.Errorf("%w: player ID %s already exists", ErrPlayerExists, id)
}
