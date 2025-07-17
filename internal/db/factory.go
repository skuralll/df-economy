package db

import (
	"github.com/skuralll/dfeconomy/db"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

func NewDB(dbType db.DBType, dsn string) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch dbType {
	case db.SQLite:
		dialector = sqlite.Dialector{DriverName: "sqlite", DSN: dsn}
	case db.MySQL:
		dialector = mysql.Open(dsn)
	case db.Postgres:
		dialector = postgres.Open(dsn)
	default:
		return nil, NewDatabaseError("create dialector", "unsupported database type: "+string(dbType))
	}

	return gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}
