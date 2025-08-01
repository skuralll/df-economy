package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/pelletier/go-toml"
	"github.com/skuralll/df-permission/permission"
	"github.com/skuralll/dfeconomy/dragonfly/commands"
	"github.com/skuralll/dfeconomy/economy/config"
	"github.com/skuralll/dfeconomy/economy/service"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	chat.Global.Subscribe(chat.StdoutSubscriber{})
	conf, err := readConfig(slog.Default())
	if err != nil {
		panic(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	cfg := config.Config{
		DBType:         "sqlite",
		DBDSN:          "./economy.db",
		DefaultBalance: 100.0,
		EnableSetCmd:   false, // Disable set command by default for security
	}

	pMgr := permission.NewManager()
	if err := pMgr.Save(); err != nil {
		slog.Error("Failed to save permission manager", "error", err)
		os.Exit(1)
	}
	slog.Info("Permission manager initialized")

	svc, cleanup, err := service.NewEconomyService(cfg, pMgr)
	if err != nil {
		slog.Error("Failed to create economy service", "error", err)
		os.Exit(1)
	}
	defer cleanup()
	commands.RegisterCommands(svc, cfg)

	srv.Listen()
	for p := range srv.Accept() {
		_ = p
		svc.RegisterUser(context.Background(), p.UUID(), p.Name())
	}
}

// readConfig reads the configuration from the config.toml file, or creates the
// file if it does not yet exist.
func readConfig(log *slog.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return c.Config(log)
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}
