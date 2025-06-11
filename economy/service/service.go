package service

import (
	"database/sql"

	"github.com/skuralll/dfeconomy/internal/db"
)

type EconomyService struct {
	db db.DB
}

// Get new EconomyService instance
func NewEconomyService(dbsql *sql.DB) *EconomyService {
	dbInstance := db.NewSQLite(dbsql) // TODO: Support multiple databases
	return &EconomyService{dbInstance}
}
