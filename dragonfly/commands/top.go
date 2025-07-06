package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"

	dfErrors "github.com/skuralll/dfeconomy/errors"
)

// /economy top <page>

const itemCount int = 10 // Number of items per page

type EconomyTopCommand struct {
	*BaseCommand
	SubCmd cmd.SubCommand `cmd:"top" help:"Show the top players by balance."`
	Page   int            `cmd:"page" help:"The page to show."`
}

func (e EconomyTopCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Execute as a player")
		return
	}

	// Provide immediate feedback
	o.Printf("Loading top balances...")

	go func() {
		// create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// get top entries
		entries, err := e.svc.GetTopBalances(ctx, e.Page, itemCount)
		if err != nil {
			switch {
			case errors.Is(err, dfErrors.ErrValueMustBeAtLeastOne):
				p.Message("§c[Error] Size must be at least 1")
			case errors.Is(err, dfErrors.ErrPageNotFound):
				p.Message("§c[Error] Page not found")
			case errors.Is(err, context.DeadlineExceeded):
				p.Message("§c[Error] Request timeout")
			default:
				p.Message("§c[Error] Failed to get top balances by internal error")
			}
			return
		}
		// success - display results
		p.Message(fmt.Sprintf("§a[Top Balances - Page %d]", e.Page))
		for i, entry := range entries {
			p.Message(fmt.Sprintf("#%d %s: %.2f", (e.Page-1)*itemCount+i+1, entry.Name, entry.Money))
		}
	}()
}

// Validation
var _ cmd.Runnable = (*EconomyTopCommand)(nil)
