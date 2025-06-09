package sqlite

import (
	"database/sql"

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
			id TEXT PRIMARY KEY,
			balance REAL NOT NULL DEFAULT 0
		);
	`)
	return err

}

func (s *svc) Balance(id uuid.UUID) (float64, error) {
	// TODO
	return 0, nil
}

func (s *svc) Set(id uuid.UUID, amount float64) error {
	// TODO
	return nil
}

func (s *svc) Top(page, size int) ([]economy.Entry, error) {
	// TODO
	return nil, nil
}
