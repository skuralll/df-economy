package commands

import (
	"context"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

// /economy set <target> <amount>

type EconomySetCommand struct {
	*BaseCommand
	SubCmd   cmd.SubCommand `cmd:"set" help:"Set the balance of a player."`
	Username string         `cmd:"username"`
	Amount   float64        `cmd:"amount"`
}

func (e EconomySetCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Execute as a player")
		return
	}
	// validate amount
	if e.Amount < 0 {
		o.Error("Amount must be at least 0")
		return
	}
	// get target uuid
	var uid uuid.UUID
	if e.Username == p.Name() {
		uid = p.UUID()
	} else {
		var err error
		uid, err = e.svc.GetUUIDByName(context.Background(), e.Username)
		if err != nil {
			o.Error("Player not found: " + e.Username)
			return
		}
	}
	// set balance
	err := e.svc.SetBalance(context.Background(), uid, e.Username, float64(e.Amount))
	if err != nil {
		o.Error("Failed to set balance: " + err.Error())
		return
	}
	o.Printf("Set balance of %s to %.2f", e.Username, e.Amount)
}

// Validation
var _ cmd.Runnable = (*EconomySetCommand)(nil)