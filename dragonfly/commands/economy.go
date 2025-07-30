package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
)

// /economy

type EconomyCommand struct {
	*BaseCommand
}

func (c EconomyCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	o.Printf("§6=== Economy Commands ===")
	o.Printf("§a/economy balance [username]§r - Display balance of yourself or another player")
	o.Printf("§a/economy pay <username> <amount>§r - Pay money to another player")
	o.Printf("§a/economy set <username> <amount>§r - Set a player's balance (Admin)")
	o.Printf("§a/economy top <page>§r - Show top players by balance")
}

// Validation
var _ cmd.Runnable = (*EconomyCommand)(nil)
var _ cmd.Allower = (*EconomyCommand)(nil)
