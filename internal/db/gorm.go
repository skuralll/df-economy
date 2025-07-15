package db

import (
	"context"
	"errors"
	"log/slog"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/skuralll/dfeconomy/economy"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

type DBGorm struct {
	db *gorm.DB
}

func NewDBGorm() (*DBGorm, func(), error) {
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        "test.db",
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		slog.Error("failed to open database", "error", err)
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to get sql.DB", "error", err)
		return nil, nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		slog.Error("database ping failed", "error", err)
		return nil, nil, err
	}

	cleanup := func() {
		// avoid null pointer
		if sqlDB, err := db.DB(); err != nil {
			slog.Error("failed to get sql.DB for cleanup", "error", err)
		} else if err := sqlDB.Close(); err != nil {
			slog.Error("failed to close database connection", "error", err)
		}
	}

	// Migrate the schema
	if err := migrateSchema(db); err != nil {
		return nil, nil, err
	}

	return &DBGorm{db}, cleanup, nil
}

// MigrateSchema migrates the database schema for the Account model.
func migrateSchema(db *gorm.DB) error {
	if err := db.AutoMigrate(&Account{}); err != nil {
		slog.Error("failed to migrate schema", "error", err)
		return err
	}
	return nil
}

func (d *DBGorm) Balance(ctx context.Context, id uuid.UUID) (float64, error) {
	// Basic data integrity checks
	if id == uuid.Nil {
		return 0, NewValidationError("uuid", "cannot be nil")
	}

	var account Account
	err := d.db.WithContext(ctx).Where("uuid = ?", id).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, NewNotFoundError("player")
		}
		return 0, NewDatabaseError("balance query", err.Error())
	}
	return account.Balance, nil
}

func (d *DBGorm) GetUUIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	// Basic data integrity checks
	if strings.TrimSpace(name) == "" {
		return uuid.Nil, NewValidationError("name", "cannot be empty")
	}

	var uStr string
	err := d.db.WithContext(ctx).Model(&Account{}).Select("uuid").Where("name = ?", name).Scan(&uStr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uuid.Nil, NewNotFoundError("player")
		}
		return uuid.Nil, NewDatabaseError("uuid query", err.Error())
	}
	// convert string to uuid
	uId, err := uuid.Parse(uStr)
	if err != nil {
		return uuid.Nil, NewDatabaseError("uuid parse", err.Error())
	}
	return uId, nil
}

func (d *DBGorm) Set(ctx context.Context, id uuid.UUID, name string, balance float64) error {
	// Basic data integrity checks
	if id == uuid.Nil {
		return NewValidationError("uuid", "cannot be nil")
	}
	if strings.TrimSpace(name) == "" {
		return NewValidationError("name", "cannot be empty")
	}
	if math.IsNaN(balance) || math.IsInf(balance, 0) {
		return NewValidationError("balance", "must be a valid number")
	}

	result := d.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.AssignmentColumns([]string{"balance", "name"}),
	}).Create(&Account{
		UUID:    id.String(),
		Name:    name,
		Balance: balance,
	})

	if result.Error != nil {
		return NewDatabaseError("balance update", result.Error.Error())
	}

	return nil
}

func (d *DBGorm) Top(ctx context.Context, page int, size int) ([]economy.EconomyEntry, error) {
	// Basic data integrity checks
	if page <= 0 {
		return nil, NewValidationError("page", "must be greater than 0")
	}
	if size <= 0 {
		return nil, NewValidationError("size", "must be greater than 0")
	}

	offset := (page - 1) * size

	// Fetch top accounts from the database
	var accounts []Account
	err := d.db.WithContext(ctx).Model(&Account{}).Limit(size).Offset(offset).Order("balance DESC").Find(&accounts).Error
	if err != nil {
		return nil, NewDatabaseError("top query", err.Error())
	}

	// Convert accounts to EconomyEntry
	var entries []economy.EconomyEntry
	for _, account := range accounts {
		u, err := uuid.Parse(string(account.UUID))
		if err != nil {
			continue // skip broken uuid
		}
		entries = append(entries, economy.EconomyEntry{
			UUID:    u,
			Name:    account.Name,
			Balance: account.Balance,
		})
	}
	return entries, nil
}

func (d *DBGorm) Transfer(ctx context.Context, fromID uuid.UUID, toID uuid.UUID, amount float64) error {
	// Basic data integrity checks
	if fromID == uuid.Nil {
		return NewValidationError("from_uuid", "cannot be nil")
	}
	if toID == uuid.Nil {
		return NewValidationError("to_uuid", "cannot be nil")
	}
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return NewValidationError("amount", "must be a valid number")
	}
	if amount <= 0 {
		return NewValidationError("amount", "must be positive")
	}

	return d.db.Transaction(func(tx *gorm.DB) error {
		// Check sender exists and get balance
		var fromAccount Account
		err := tx.Where("uuid = ?", fromID).First(&fromAccount).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NewNotFoundError("sender")
			}
			return NewDatabaseError("sender query", err.Error())
		}
		if fromAccount.Balance < amount {
			return NewInsufficientBalanceError(amount, fromAccount.Balance)
		}
		// Check receiver exists
		err = tx.Model(&Account{}).Where("uuid = ?", toID).First(&Account{}).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NewNotFoundError("receiver")
			}
			return NewDatabaseError("receiver query", err.Error())
		}
		// Deduct from sender
		result := tx.Model(&Account{}).Where("uuid = ?", fromID).Update("balance", gorm.Expr("balance - ?", amount))
		if result.Error != nil {
			return NewDatabaseError("sender balance update", result.Error.Error())
		}
		// Add to receiver
		result = tx.Model(&Account{}).Where("uuid = ?", toID).Update("balance", gorm.Expr("balance + ?", amount))
		if result.Error != nil {
			return NewDatabaseError("receiver balance update", result.Error.Error())
		}
		// Return nil to indicate success
		return nil
	})
}

// Implementation completeness checks
var _ DB = (*DBGorm)(nil)
