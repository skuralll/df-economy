package db

import (
	"gorm.io/gorm"
)

// Accounts represents a user account in the database.
type Account struct {
	gorm.Model
	UUID    string  `gorm:"type:char(36);uniqueIndex;not null"`
	Name    string  `gorm:"type:varchar(16);not null"`
	Balance float64 `gorm:"type:real;not null;default:0"`
}
