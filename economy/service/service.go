package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/skuralll/dfeconomy/economy/config"
	"github.com/skuralll/dfeconomy/errors"
	"github.com/skuralll/dfeconomy/internal/db"
	"github.com/skuralll/dfeconomy/models"
)

type EconomyService struct {
	db  db.DB
	cfg config.Config
}

// Get new EconomyService instance
func NewEconomyService(cfg config.Config) (*EconomyService, func(), error) {
	dbInstance, cleanup, err := db.NewSQLiteFromConfig(&db.SQLiteConfig{Path: cfg.DBPath}) // TODO: Support multiple databases
	if err != nil {
		return nil, nil, err
	}
	return &EconomyService{dbInstance, cfg}, cleanup, nil
}

// TODO: Move validation logic in db to the service. The db should only operate the database based on the received values.

// Register a new user
func (svc *EconomyService) RegisterUser(ctx context.Context, id uuid.UUID, name string) (bool, error) {
	// Check if user already exists
	_, err := svc.db.Balance(ctx, id)
	if err == nil {
		// User already exists
		return false, nil
	}
	// Register new user
	err = svc.db.Set(ctx, id, name, svc.cfg.DefaultAmount)
	if err != nil {
		return false, err
	}
	slog.Info("New user registered", "id", id, "name", name)
	return true, nil
}

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

// Transfer balance
func (svc *EconomyService) TransferBalance(ctx context.Context, fromID, toID uuid.UUID, amount float64) error {
	if fromID == toID {
		return errors.ErrCannotTargetSelf
	}
	err := svc.db.Transfer(ctx, fromID, toID, amount)
	return err
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
	// error handle
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.ErrPageNotFound
	}
	return list, nil
}

// GetUUIDByName retrieves the UUID of a player by their name.
func (svc *EconomyService) GetUUIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	uid, err := svc.db.GetUUIDByName(ctx, name)
	if err != nil {
		return uuid.Nil, errors.ErrUnknownPlayer
	}

	return uid, nil
}
