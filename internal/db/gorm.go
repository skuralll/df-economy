package db

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/skuralll/dfeconomy/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

type DBGorm struct {
	db *gorm.DB
}

func NewDBGorm() (*DBGorm, func(), error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
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

// Balance implements DB.
func (d *DBGorm) Balance(ctx context.Context, id uuid.UUID) (float64, error) {
	panic("unimplemented")
}

// GetUUIDByName implements DB.
func (d *DBGorm) GetUUIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	panic("unimplemented")
}

// Set implements DB.
func (d *DBGorm) Set(ctx context.Context, id uuid.UUID, name string, amount float64) error {
	panic("unimplemented")
}

// Top implements DB.
func (d *DBGorm) Top(ctx context.Context, page int, size int) ([]models.EconomyEntry, error) {
	panic("unimplemented")
}

// Transfer implements DB.
func (d *DBGorm) Transfer(ctx context.Context, fromID uuid.UUID, toID uuid.UUID, amount float64) error {
	panic("unimplemented")
}

// Implementation completeness checks
var _ DB = (*DBGorm)(nil)
