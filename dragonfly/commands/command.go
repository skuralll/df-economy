package commands

import "github.com/df-mc/dragonfly/server/cmd"

func RegisterCommands() {
	cmd.Register(cmd.New("economy", "Displays economy-related information.", nil, EconomyCommand{}))
}
