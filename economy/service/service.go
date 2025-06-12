package service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/skuralll/dfeconomy/internal/db"
)

type EconomyService struct {
	db db.DB
}

// Get new EconomyService instance
func NewEconomyService(dbsql *sql.DB) *EconomyService {
	dbInstance := db.NewSQLite(dbsql) // TODO: Support multiple databases
	return &EconomyService{dbInstance}
}

// Set balance
// Set(ctx context.Context, id uuid.UUID, name *string, amount float64) error
// Get balance ranking
// Top(ctx context.Context, page, size int) ([]EconomyEntry, error)

// TODO: Move validation logic in db to the service. The db should only operate the database based on the received values.

// Get balance
func (svc *EconomyService) GetBalance(ctx context.Context, id uuid.UUID) (float64, error) {
	amount, err := svc.db.Balance(ctx, id)
	if err != nil {
		return 0, err
	}
	return amount, nil
}
