package config

type Config struct {
	DBPath         string  `toml:"db_path"`        // Path to the database file
	DefaultBalance float64 `toml:"default_amount"` // Default amount of money for new users
}
