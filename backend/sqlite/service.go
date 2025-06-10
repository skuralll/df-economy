package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/skuralll/dfeconomy/economy"
)

type svc struct {
	db *sql.DB
}

func New(db *sql.DB) economy.Economy {
	if db == nil {
		panic("db cannot be nil")
	}
	if err := initSchema(db); err != nil {
		panic(err)
	}
	return &svc{db: db}
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

func (s *svc) Balance(ctx context.Context, id uuid.UUID) (float64, error) {
	var amount float64
	err := s.db.QueryRowContext(ctx, "SELECT money FROM balances WHERE uuid = ?", id.String()).Scan(&amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, economy.ErrUnknownPlayer
		}
		return 0, err
	}
	return amount, nil
}

func (s *svc) Set(ctx context.Context, id uuid.UUID, name *string, amount float64) error {
	if amount < 0 {
		return economy.ErrNegativeAmount
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO balances (uuid, name ,money) VALUES (?,?,?)
		ON CONFLICT (uuid) DO UPDATE 
			SET money = excluded.money,
					name = COALESCE(excluded.name, balances.name)
	`, id.String(), name, amount)
	return err
}

func (s *svc) Top(
	ctx context.Context,
	page, size int, // page 1-based, size > 0
) ([]economy.Entry, error) {

	if size <= 0 {
		return nil, errors.New("size must be positive")
	}
	if page <= 0 {
		page = 1
	}
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

	var list []economy.Entry
	for rows.Next() {
		var (
			uStr  string
			name  sql.NullString
			money int64
		)
		if err := rows.Scan(&uStr, &name, &money); err != nil {
			return nil, err
		}
		u, err := uuid.Parse(uStr)
		if err != nil { // skip broken uuid
			continue
		}
		list = append(list, economy.Entry{
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
