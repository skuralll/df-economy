package db

import (
	"log/slog"

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

	return &DBGorm{db}, cleanup, nil
}
