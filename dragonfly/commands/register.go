package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/skuralll/dfeconomy/economy/config"
	"github.com/skuralll/dfeconomy/economy/service"
)

func RegisterCommands(svc *service.EconomyService, cfg config.Config) {
	baseCmd := &BaseCommand{svc: svc}
	
	// Base commands that are always available
	subCommands := []cmd.Runnable{
		&EconomyBalanceCommand{BaseCommand: baseCmd},
		&EconomyTopCommand{BaseCommand: baseCmd},
		&EconomyPayCommand{BaseCommand: baseCmd},
		&EconomyCommand{baseCmd},
	}
	
	// Conditionally add set command
	if cfg.EnableSetCmd {
		subCommands = append(subCommands, &EconomySetCommand{BaseCommand: baseCmd})
	}
	
	cmd.Register(cmd.New("economy", "Displays economy-related information.", nil, subCommands...))
}