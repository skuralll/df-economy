package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// /economy

type EconomyCommand struct {
	*BaseCommand
}

func (c EconomyCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	o.Printf("Executed: /economy")
}

// /economy balance <target>

type EconomyBalanceCommand struct {
	*BaseCommand
	SubCmd   cmd.SubCommand       `cmd:"balance" help:"Displays the balance of a player."`
	Username cmd.Optional[string] `cmd:"username"`
}

func (e EconomyBalanceCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Execute as a player")
	}

	o.Printf("Executed: /economy balance %s", e.Username.LoadOr(p.Name()))
}

// Validation
var _ cmd.Runnable = (*EconomyCommand)(nil)
var _ cmd.Runnable = (*EconomyBalanceCommand)(nil)
