package dfeconomy

import (
	"log/slog"

	"github.com/df-mc/dragonfly/server"
)

// plugin instance
type DfEconomyPlugin struct{}

// creates a new plugin instance
func NewDfEconomyPlugin() *DfEconomyPlugin {
	return &DfEconomyPlugin{}
}

func (p *DfEconomyPlugin) Enable(srv *server.Server) error {
	slog.Info("DfEconomy plugin enabled")
	return nil
}

func (p *DfEconomyPlugin) Disable() error {
	slog.Info("DfEconomy plugin enabled")
	return nil
}
