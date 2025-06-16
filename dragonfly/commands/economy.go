package commands

import (
	"context"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

// /economy

type EconomyCommand struct {
	*BaseCommand
}

func (c EconomyCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	o.Printf("help: TODO")
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
	// get target uuid
	tn := e.Username.LoadOr(p.Name())
	var uid uuid.UUID
	if tn == p.Name() {
		uid = p.UUID()
	} else {
		// get uuid by name
		var err error
		uid, err = e.svc.GetUUIDByName(context.Background(), tn)
		if err != nil {
			o.Error("Player not found: " + tn)
			return
		}
	}
	// get balance
	amount, err := e.svc.GetBalance(context.Background(), uid)
	if err != nil {
		o.Error("Failed to get balance")
		return
	}
	o.Printf("Balance of %s: %.2f", tn, amount)
}

// /economy set <target> <amount>

type EconomySetCommand struct {
	*BaseCommand
	SubCmd   cmd.SubCommand `cmd:"set" help:"Set the balance of a player."`
	Username string         `cmd:"username"`
	Amount   float32        `cmd:"amount"`
}

func (e EconomySetCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {

}

// Validation
var _ cmd.Runnable = (*EconomyCommand)(nil)
var _ cmd.Runnable = (*EconomyBalanceCommand)(nil)
var _ cmd.Runnable = (*EconomySetCommand)(nil)
