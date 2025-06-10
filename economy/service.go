package economy

import (
	"context"

	"github.com/google/uuid"
)

// Economy defines the interface for managing user balances and transactions
// within an economic system. It provides methods to query balances, set balances,
// and retrieve a leaderboard.
//
// Methods:
//   - Balance(id uuid.UUID): Returns the current balance for the specified user.
//   - Set(id uuid.UUID, amount float64): Sets the user's balance to the specified amount.
//   - Top(page, size int): Retrieves a paginated list of top entries based on balance.
type Economy interface {
	Balance(ctx context.Context, id uuid.UUID) (float64, error)
	Set(ctx context.Context, id uuid.UUID, name *string, amount float64) error
	Top(ctx context.Context, page, size int) ([]Entry, error)
}

// Entry represents a single row in the balance leaderboard.
//
// UUID   is the unique identifier of the player.
// Name   is the most recently known display name of the player.
// Money  is the player’s current balance expressed in the smallest unit.
//
// The slice returned by Economy.Top is ordered by Money in descending order.
type Entry struct {
	UUID  uuid.UUID // Player’s UUID
	Name  string    // Display name
	Money float64   // Balance
}
