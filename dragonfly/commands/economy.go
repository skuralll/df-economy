package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
)

type EconomyCommand struct{}

// Run implements cmd.Runnable.
func (c EconomyCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	o.Printf("実行されました: /economy")
}

// Validate implements cmd.Runnable.
var _ cmd.Runnable = (*EconomyCommand)(nil)
