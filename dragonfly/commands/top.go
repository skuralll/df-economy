package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"

	"github.com/skuralll/dfeconomy/economy/service"
)

// /economy top <page>

const itemCount int = 10 // Number of items per page

type EconomyTopCommand struct {
	*BaseCommand
	SubCmd cmd.SubCommand `cmd:"top" help:"Show the top players by balance."`
	Page   int            `cmd:"page" help:"The page to show."`
}

func (e *EconomyTopCommand) Allow(src cmd.Source) bool {
	return e.CheckPermission(src, "economy.command.top")
}

func (e EconomyTopCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := e.ValidatePlayerSource(src, o)
	if !ok {
		return
	}

	// Provide immediate feedback
	o.Printf("Loading top balances...")

	e.ExecuteAsync(p, func(ctx context.Context) {
		// get top entries
		entries, err := e.svc.GetTopBalances(ctx, e.Page, itemCount)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrValidation):
				p.Message("§c[Error] Invalid input: " + err.Error())
			case errors.Is(err, context.DeadlineExceeded):
				p.Message("§c[Error] Request timeout")
			case errors.Is(err, service.ErrInternalError):
				p.Message("§c[Error] Failed to get top balances by internal error")
			default:
				p.Message("§c[Error] Failed to get top balances by internal error")
			}
			return
		}
		// success - display results
		p.Message(fmt.Sprintf("§a[Top Balances - Page %d]", e.Page))
		for i, entry := range entries {
			p.Message(fmt.Sprintf("#%d %s: %.2f", (e.Page-1)*itemCount+i+1, entry.Name, entry.Balance))
		}
	})
}

// Validation
var _ cmd.Runnable = (*EconomyTopCommand)(nil)
var _ cmd.Allower = (*EconomyTopCommand)(nil)
