package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/skuralll/dfeconomy/errors"
	"github.com/skuralll/dfeconomy/internal/db"
	"github.com/skuralll/dfeconomy/models"
)

type EconomyService struct {
	db db.DB
}

// Get new EconomyService instance
func NewEconomyService() (*EconomyService, func(), error) {
	dbInstance, cleanup, err := db.NewSQLiteFromConfig(&db.SQLiteConfig{Path: "./foo.db"}) // TODO: Support multiple databases
	if err != nil {
		return nil, nil, err
	}
	return &EconomyService{dbInstance}, cleanup, nil
}

// TODO: Move validation logic in db to the service. The db should only operate the database based on the received values.

// Get balance
func (svc *EconomyService) GetBalance(ctx context.Context, id uuid.UUID) (float64, error) {
	amount, err := svc.db.Balance(ctx, id)
	if err != nil {
		return 0, err
	}
	return amount, nil
}

// Set balance
func (svc *EconomyService) SetBalance(ctx context.Context, id uuid.UUID, name string, amount float64) error {
	if amount < 0 {
		return errors.ErrNegativeAmount
	}
	result := svc.db.Set(ctx, id, name, amount)
	return result
}

// Get balance ranking
func (svc *EconomyService) GetTopBalances(ctx context.Context, page, size int) ([]models.EconomyEntry, error) {
	// validation
	if size <= 0 {
		return nil, errors.ErrValueMustBeAtLeastOne
	}
	if page <= 0 {
		return nil, errors.ErrPageNotFound
	}
	// get result
	list, err := svc.db.Top(ctx, page, size)
	if err != nil {
		return nil, err
	}
	return list, nil
}
