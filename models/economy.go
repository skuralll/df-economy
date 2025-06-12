package models

import "github.com/google/uuid"

type EconomyEntry struct {
	UUID  uuid.UUID // Player’s UUID
	Name  string    // Display name
	Money float64   // Balance
}
