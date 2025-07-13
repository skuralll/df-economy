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

func (d *DBGorm) Balance(ctx context.Context, id uuid.UUID) (float64, error) {
	var balance float64
	err := d.db.Model(&Account{}).Select("balance").Where("uuid = ?", id).Scan(&balance).Error
	if err != nil {
		slog.Error("failed to get balance", "uuid", id, "error", err)
		return 0, err
	}
	return balance, nil
}

// GetUUIDByName implements DB.
func (d *DBGorm) GetUUIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	var uStr string
	err := d.db.Model(&Account{}).Select("uuid").Where("name = ?", name).Scan(&uStr).Error
	if err != nil {
		slog.Error("failed to get uuid by name", "name", name, "error", err)
		return uuid.Nil, err
	}
	// convert string to uuid
	uId, err := uuid.Parse(uStr)
	if err != nil {
		return uuid.Nil, err
	}
	return uId, nil
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
