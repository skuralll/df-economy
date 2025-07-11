package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/skuralll/dfeconomy/economy/service"
)

func RegisterCommands(svc *service.EconomyService) {
	baseCmd := &BaseCommand{svc: svc}
	// Fill in the required fields for EconomyBalanceCommand as needed
	cmd.Register(cmd.New("economy", "Displays economy-related information.", nil,
		&EconomyBalanceCommand{BaseCommand: baseCmd},
		&EconomySetCommand{BaseCommand: baseCmd},
		&EconomyTopCommand{BaseCommand: baseCmd},
		&EconomyPayCommand{BaseCommand: baseCmd},
		&EconomyCommand{baseCmd},
	))
}