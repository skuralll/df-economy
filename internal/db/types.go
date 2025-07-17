package db

type DBType string

const (
	SQLite   DBType = "sqlite"
	MySQL    DBType = "mysql"
	Postgres DBType = "postgres"
)
