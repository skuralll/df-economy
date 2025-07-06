package commands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
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
	p, ok := e.ValidatePlayerSource(src, o)
	if !ok {
		return
	}

	// Provide immediate feedback
	o.Printf("Processing payment...")

	go func() {
		// create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// get target uuid
		tuid, err := e.svc.GetUUIDByName(ctx, e.Username)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				p.Message("§c[Error] Request timeout")
			} else {
				p.Message("§c[Error] Player not found: " + e.Username)
			}
			return
		}
		err = e.svc.TransferBalance(ctx, p.UUID(), tuid, e.Amount)
		if err != nil {
			switch {
			case errors.Is(err, dfErrors.ErrValueMustBeAtLeastOne):
				p.Message("§c[Error] Amount must be at least 1")
			case errors.Is(err, dfErrors.ErrInsufficientFunds):
				p.Message("§c[Error] Insufficient funds")
			case errors.Is(err, dfErrors.ErrUnknownPlayer):
				p.Message("§c[Error] Target player not found: " + e.Username)
			case errors.Is(err, context.DeadlineExceeded):
				p.Message("§c[Error] Request timeout")
			case errors.Is(err, dfErrors.ErrCannotTargetSelf):
				p.Message("§c[Error] Cannot pay yourself")
			case errors.Is(err, dfErrors.ErrNegativeAmount):
				p.Message("§c[Error] Amount must be positive")
			default:
				p.Message("§c[Error] Failed to pay by internal error")
				slog.Error("Failed to pay", "error", err, "from", p.Name(), "to", e.Username, "amount", e.Amount)
			}
			return
		}
		// success
		p.Message(fmt.Sprintf("§a[Success] You paid %.2f to %s", e.Amount, e.Username))
	}()
}

// Validation
var _ cmd.Runnable = (*EconomyPayCommand)(nil)
