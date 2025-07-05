package commands

import (
	"github.com/skuralll/dfeconomy/economy/service"
)

type BaseCommand struct {
	svc *service.EconomyService
}