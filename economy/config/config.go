package config

type Config struct {
	DBType         string  `toml:"db_type"`         // Database type: sqlite, mysql, postgres
	DBDSN          string  `toml:"db_dsn"`          // Path to the database file
	DefaultBalance float64 `toml:"default_balance"` // Default amount of money for new users
	EnableSetCmd   bool    `toml:"enable_set_cmd"`  // Enable /economy set command
}
