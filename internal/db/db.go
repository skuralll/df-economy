package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/skuralll/dfeconomy/models"
)

type DB interface {
	// Get balance
	Balance(ctx context.Context, id uuid.UUID) (float64, error)
	// Set balance
	Set(ctx context.Context, id uuid.UUID, name string, amount float64) error
	// Get balance ranking
	Top(ctx context.Context, page, size int) ([]models.EconomyEntry, error)
	// Get uuid by name
	GetUUIDByName(ctx context.Context, name string) (uuid.UUID, error)
}
