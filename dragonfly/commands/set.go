package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
)

// /economy set <target> <amount>

type EconomySetCommand struct {
	*BaseCommand
	SubCmd   cmd.SubCommand `cmd:"set" help:"Set the balance of a player."`
	Username string         `cmd:"username"`
	Amount   float64        `cmd:"amount"`
}

func (e EconomySetCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := e.ValidatePlayerSource(src, o)
	if !ok {
		return
	}
	// validate amount
	if e.Amount < 0 {
		o.Error("Amount must be at least 0")
		return
	}

	// Provide immediate feedback
	o.Printf("Processing balance update...")

	e.ExecuteAsync(p, func(ctx context.Context) {
		// get target uuid
		tuid, err := e.GetUUIDByName(ctx, p, e.Username)
		if err != nil {
			return
		}
		// set balance
		err = e.svc.SetBalance(ctx, tuid, e.Username, float64(e.Amount))
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				p.Message("§c[Error] Request timeout")
			} else {
				p.Message("§c[Error] Failed to set balance: " + err.Error())
			}
			return
		}
		// success
		p.Message(fmt.Sprintf("§a[Success] Set balance of %s to %.2f", e.Username, e.Amount))
	})
}

// Validation
var _ cmd.Runnable = (*EconomySetCommand)(nil)
var _ cmd.Allower = (*EconomySetCommand)(nil)
