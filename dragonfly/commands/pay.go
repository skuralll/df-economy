package commands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"

	"github.com/skuralll/dfeconomy/economy/service"
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

	e.ExecuteAsync(p, func(ctx context.Context) {
		// get target uuid
		tuid, err := e.GetUUIDByName(ctx, p, e.Username)
		if err != nil {
			return
		}
		// transfer balance
		err = e.svc.TransferBalance(ctx, p.UUID(), tuid, e.Amount)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrValidation):
				p.Message("§c[Error] Invalid input: " + err.Error())
			case errors.Is(err, service.ErrUnknownPlayer):
				p.Message("§c[Error] Target player not found: " + e.Username)
			case errors.Is(err, context.DeadlineExceeded):
				p.Message("§c[Error] Request timeout")
			case errors.Is(err, service.ErrInternalError):
				p.Message("§c[Error] Failed to pay by internal error")
				slog.Error("Failed to pay", "error", err, "from", p.Name(), "to", e.Username, "amount", e.Amount)
			default:
				p.Message("§c[Error] Failed to pay by internal error")
				slog.Error("Failed to pay", "error", err, "from", p.Name(), "to", e.Username, "amount", e.Amount)
			}
			return
		}
		// success
		p.Message(fmt.Sprintf("§a[Success] You paid %.2f to %s", e.Amount, e.Username))
	})
}

// Validation
var _ cmd.Runnable = (*EconomyPayCommand)(nil)
var _ cmd.Allower = (*EconomyPayCommand)(nil)
