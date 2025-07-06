package commands

import (
	"context"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/skuralll/dfeconomy/economy/service"
)

const (
	DefaultCommandTimeout = 5 * time.Second
)

type BaseCommand struct {
	svc *service.EconomyService
}

// ValidatePlayerSource validates that the source is a player and outputs error if not.
// This method intentionally combines validation and error output for code brevity.
func (b *BaseCommand) ValidatePlayerSource(src cmd.Source, o *cmd.Output) (*player.Player, bool) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Execute as a player")
	}
	return p, ok
}

// Create Context with Timeout creates a context with a 5-second timeout.
func (b *BaseCommand) CreateContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultCommandTimeout)
}
