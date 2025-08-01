package economy

import "github.com/google/uuid"

type EconomyEntry struct {
	UUID    uuid.UUID // Player's UUID
	Name    string    // Display name
	Balance float64   // Balance
}
