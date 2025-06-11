package db

import (
	"context"

	"github.com/google/uuid"
)

type EconomyEntry struct {
	UUID  uuid.UUID // Player’s UUID
	Name  string    // Display name
	Money float64   // Balance
}

type DB interface {
	// Get balance
	Balance(ctx context.Context, id uuid.UUID) (float64, error)
	// Set balance
	Set(ctx context.Context, id uuid.UUID, name *string, amount float64) error
	// Get balance ranking
	Top(ctx context.Context, page, size int) ([]EconomyEntry, error)
}
