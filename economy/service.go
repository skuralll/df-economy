package economy

import "github.com/google/uuid"

// Economy defines the interface for managing user balances and transactions
// within an economic system. It provides methods to query balances, deposit and
// withdraw funds, transfer amounts between accounts, and retrieve a leaderboard.
//
// Methods:
//   - Balance(id uuid.UUID): Returns the current balance for the specified user.
//   - Deposit(id uuid.UUID, amount float64): Adds the specified amount to the user's balance.
//   - Withdraw(id uuid.UUID, amount float64): Deducts the specified amount from the user's balance.
//   - Transfer(from, to uuid.UUID, amount float64): Transfers the specified amount from one user to another.
//   - Top(page, size int): Retrieves a paginated list of top entries based on balance.
type Economy interface {
	Balance(id uuid.UUID) (float64, error)
	Deposit(id uuid.UUID, amount float64) error
	Withdraw(id uuid.UUID, amount float64) error
	Transfer(from, to uuid.UUID, amount float64) error
	Top(page, size int) ([]Entry, error)
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
	Money int64     // Balance
}
