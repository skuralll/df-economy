package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/skuralll/dfeconomy/economy"
	"github.com/skuralll/dfeconomy/economy/config"
	"github.com/skuralll/dfeconomy/internal/db"
)

type EconomyService struct {
	db  db.DB
	cfg config.Config
}

// Get new EconomyService instance
func NewEconomyService(cfg config.Config) (*EconomyService, func(), error) {
	dbInstance, cleanup, err := db.NewDBGorm() // TODO: Support multiple databases
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
		return false, NewPlayerExistsError(id.String())
	}
	// Register new user
	err = svc.db.Set(ctx, id, name, svc.cfg.DefaultBalance)
	if err != nil {
		if errors.Is(err, db.ErrValidation) {
			return false, NewValidationError("user data", err.Error())
		}
		return false, NewInternalError("user registration", err.Error())
	}
	slog.Info("New user registered", "id", id, "name", name)
	return true, nil
}

// Get balance
func (svc *EconomyService) GetBalance(ctx context.Context, id uuid.UUID) (float64, error) {
	amount, err := svc.db.Balance(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return 0, NewUnknownPlayerError(id.String())
		}
		return 0, NewInternalError("balance query", err.Error())
	}
	return amount, nil
}

// Set balance
func (svc *EconomyService) SetBalance(ctx context.Context, id uuid.UUID, name string, amount float64) error {
	if amount < 0 {
		return NewValidationError("amount", "must be positive")
	}
	err := svc.db.Set(ctx, id, name, amount)
	if err != nil {
		if errors.Is(err, db.ErrValidation) {
			return NewValidationError("balance data", err.Error())
		}
		return NewInternalError("balance update", err.Error())
	}
	return nil
}

// Transfer balance
func (svc *EconomyService) TransferBalance(ctx context.Context, fromID, toID uuid.UUID, amount float64) error {
	if fromID == toID {
		return NewValidationError("target", "cannot target yourself")
	}
	err := svc.db.Transfer(ctx, fromID, toID, amount)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return NewUnknownPlayerError("player in transfer")
		}
		if errors.Is(err, db.ErrInsufficientBalance) {
			return NewValidationError("balance", "insufficient funds")
		}
		if errors.Is(err, db.ErrValidation) {
			return NewValidationError("transfer data", err.Error())
		}
		return NewInternalError("transfer", err.Error())
	}
	return nil
}

// Get balance ranking
func (svc *EconomyService) GetTopBalances(ctx context.Context, page, size int) ([]economy.EconomyEntry, error) {
	// validation
	if size <= 0 {
		return nil, NewValidationError("size", "must be at least 1")
	}
	if page <= 0 {
		return nil, NewValidationError("page", "must be at least 1")
	}
	// get result
	list, err := svc.db.Top(ctx, page, size)
	// error handle
	if err != nil {
		if errors.Is(err, db.ErrValidation) {
			return nil, NewValidationError("pagination", err.Error())
		}
		return nil, NewInternalError("top query", err.Error())
	}
	if len(list) == 0 {
		return nil, NewValidationError("page", "not found")
	}
	return list, nil
}

// GetUUIDByName retrieves the UUID of a player by their name.
func (svc *EconomyService) GetUUIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	uid, err := svc.db.GetUUIDByName(ctx, name)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return uuid.Nil, NewUnknownPlayerError(name)
		}
		return uuid.Nil, NewInternalError("player lookup", err.Error())
	}

	return uid, nil
}
