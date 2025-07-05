package commands

import (
	"context"
	"errors"
	"log/slog"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"

	dfErrors "github.com/skuralll/dfeconomy/errors"
)

// /economy pay <target> <amount>

type EconomyPayCommand struct {
	*BaseCommand
	SubCmd   cmd.SubCommand `cmd:"pay" help:"Pay a player."`
	Username string         `cmd:"username"`
	Amount   float64        `cmd:"amount"`
}

func (e EconomyPayCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Execute as a player")
		return
	}
	// get target uuid
	tuid, err := e.svc.GetUUIDByName(context.Background(), e.Username)
	if err != nil {
		o.Error("Player not found: " + e.Username)
		return
	}
	err = e.svc.TransferBalance(context.Background(), p.UUID(), tuid, e.Amount)
	if err != nil {
		switch {
		case errors.Is(err, dfErrors.ErrValueMustBeAtLeastOne):
			o.Error("Amount must be at least 1")
			return
		case errors.Is(err, dfErrors.ErrInsufficientFunds):
			o.Error("Insufficient funds")
			return
		case errors.Is(err, dfErrors.ErrUnknownPlayer):
			o.Error("Target player not found: " + e.Username)
			return
		default:
			o.Error("Failed to pay by internal error")
			slog.Error("Failed to pay", "error", err, "from", p.Name(), "to", e.Username, "amount", e.Amount)
			return
		}
	}
	// success
	o.Printf("You paid %.2f to %s", e.Amount, e.Username)
}

// Validation
var _ cmd.Runnable = (*EconomyPayCommand)(nil)