package dfeconomy

import (
	"log/slog"

	"github.com/df-mc/dragonfly/server"
	"github.com/skuralll/dfeconomy/dragonfly/commands"
	"github.com/skuralll/dfeconomy/economy/service"
)

// plugin instance
type DfEconomyPlugin struct {
	svc          *service.EconomyService
	cleanupFuncs []func()
}

// creates a new plugin instance
func NewDfEconomyPlugin() *DfEconomyPlugin {
	return &DfEconomyPlugin{}
}

// Enable is called when the plugin is enabled by the server.
func (p *DfEconomyPlugin) Enable(srv *server.Server) error {
	// create economy service
	svc, cleanup, err := service.NewEconomyService()
	if err != nil {
		slog.Error("Failed to create economy service", "error", err)
	}
	p.svc = svc
	p.cleanupFuncs = append(p.cleanupFuncs, cleanup)
	// register commands
	commands.RegisterCommands(p.svc)

	slog.Info("DfEconomy plugin enabled")
	return nil
}

// Disable is called when the plugin is disabled by the server.
func (p *DfEconomyPlugin) Disable() error {
	slog.Info("DfEconomy plugin disabled")
	return nil
}
