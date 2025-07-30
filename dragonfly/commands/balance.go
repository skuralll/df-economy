package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

// /economy balance <target>

type EconomyBalanceCommand struct {
	*BaseCommand
	SubCmd   cmd.SubCommand       `cmd:"balance" help:"Displays the balance of a player."`
	Username cmd.Optional[string] `cmd:"username"`
}

func (e *EconomyBalanceCommand) Allow(src cmd.Source) bool {
	p, ok := src.(*player.Player)
	if !ok {
		return false
	}
	return e.svc.Permission.HasPermission(p.UUID(), "economy.command.balance")
}

func (e EconomyBalanceCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := e.ValidatePlayerSource(src, o)
	if !ok {
		return
	}

	// Provide immediate feedback
	o.Printf("Fetching balance...")

	e.ExecuteAsync(p, func(ctx context.Context) {
		// get target uuid
		tn := e.Username.LoadOr(p.Name())
		var uid uuid.UUID
		if tn == p.Name() {
			uid = p.UUID()
		} else {
			// get uuid by name
			var err error
			uid, err = e.GetUUIDByName(ctx, p, tn)
			if err != nil {
				return
			}
		}

		// get balance
		amount, err := e.svc.GetBalance(ctx, uid)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				p.Message("§c[Error] Request timeout")
			} else {
				p.Message("§c[Error] Failed to get balance")
			}
			return
		}
		// send message
		p.Message(fmt.Sprintf("§a[Balance] %s: %.2f", tn, amount))
	})
}

// Validation
var _ cmd.Runnable = (*EconomyBalanceCommand)(nil)
var _ cmd.Allower = (*EconomyBalanceCommand)(nil)
