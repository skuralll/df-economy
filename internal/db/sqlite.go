package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	ecerrors "github.com/skuralll/dfeconomy/errors"
	"github.com/skuralll/dfeconomy/models"
)

// Implementation of DB using SQLite
type DBSQLite struct {
	db *sql.DB
}

type SQLiteConfig struct {
	Path string
}

// Implementation completeness checks
var _ DB = (*DBSQLite)(nil)

func NewSQLiteFromConfig(config *SQLiteConfig) (*DBSQLite, func(), error) {
	db, err := sql.Open("sqlite3", config.Path)
	if err != nil {
		return nil, nil, err
	}

	// test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, nil, err
	}

	// create lambda function
	sqliteDB, cleanup, err := NewSQLite(db)
	if err != nil {
		db.Close()
		return nil, nil, err
	}

	return sqliteDB, cleanup, nil
}

func NewSQLite(db *sql.DB) (*DBSQLite, func(), error) {
	if db == nil {
		return nil, nil, errors.New("db cannot be nil")
	}
	if err := initSchema(db); err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		db.Close()
	}
	return &DBSQLite{db}, cleanup, nil
}

func initSchema(db *sql.DB) error {
	// Create the necessary tables if they do not exist.
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS balances (
			uuid TEXT PRIMARY KEY,
			name TEXT,
			money REAL NOT NULL DEFAULT 0
		);
	`)
	return err
}

func (s *DBSQLite) Balance(ctx context.Context, id uuid.UUID) (float64, error) {
	var amount float64
	err := s.db.QueryRowContext(ctx, "SELECT money FROM balances WHERE uuid = ?", id.String()).Scan(&amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ecerrors.ErrUnknownPlayer
		}
		return 0, err
	}
	return amount, nil
}

func (s *DBSQLite) Set(ctx context.Context, id uuid.UUID, name string, amount float64) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO balances (uuid, name ,money) VALUES (?,?,?)
		ON CONFLICT (uuid) DO UPDATE 
			SET money = excluded.money,
					name = COALESCE(excluded.name, balances.name)
	`, id.String(), name, amount)
	return err
}

// todo:refactor
func (s *DBSQLite) Top(
	ctx context.Context,
	page, size int, // page 1-based, size > 0
) ([]models.EconomyEntry, error) {
	offset := (page - 1) * size

	rows, err := s.db.QueryContext(ctx, `
		SELECT uuid, name, money
		FROM balances
		ORDER BY money DESC
		LIMIT ? OFFSET ?
	`, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.EconomyEntry
	for rows.Next() {
		var (
			uStr  string
			name  sql.NullString
			money float64
		)
		if err := rows.Scan(&uStr, &name, &money); err != nil {
			return nil, err
		}
		u, err := uuid.Parse(uStr)
		if err != nil { // skip broken uuid
			continue
		}
		list = append(list, models.EconomyEntry{
			UUID:  u,
			Name:  name.String,
			Money: money,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}
