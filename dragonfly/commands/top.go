package commands

import (
	"context"
	"errors"

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
	_, ok := src.(*player.Player)
	if !ok {
		o.Error("Execute as a player")
		return
	}
	entries, err := e.svc.GetTopBalances(context.Background(), e.Page, itemCount)
	if err != nil {
		if errors.Is(err, dfErrors.ErrValueMustBeAtLeastOne) {
			o.Error("Size must be at least 1")
		} else if errors.Is(err, dfErrors.ErrPageNotFound) {
			o.Error("Page not found")
		} else {
			o.Error("Failed to get top balances by internal error")
		}
		return
	}
	for i, entry := range entries {
		o.Printf("#%d %s: %.2f", (e.Page-1)*itemCount+i+1, entry.Name, entry.Money)
	}
}

// Validation
var _ cmd.Runnable = (*EconomyTopCommand)(nil)