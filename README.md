# DF Economy

Economy system library for Minecraft Bedrock Edition (**df-mc/dragonfly** framework)

## Overview

DF Economy provides a currency system for Minecraft servers with player-to-player transactions, balance checking, and leaderboard features.

## In-Game Commands

| Command | Description | Example |
| --- | --- | --- |
| `/economy` | Show command help | `/economy` |
| `/economy balance [player]` | Display balance | `/economy balance` or `/economy balance Steve` |
| `/economy pay <player> <amount>` | Send money to another player | `/economy pay Steve 100` |
| `/economy set <player> <amount>` | Set player balance (admin) | `/economy set Steve 1000` |
| `/economy top <page>` | Show balance leaderboard | `/economy top 1` |

## Usage

### 1. Install Library

```bash
go get github.com/skuralll/dfeconomy
```

### 2. Server Integration

```go
package main

import (
    "github.com/df-mc/dragonfly/server"
    "github.com/skuralll/dfeconomy/dragonfly/commands"
    "github.com/skuralll/dfeconomy/economy/service"
)

func main() {
    // Create economy service
    svc, cleanup, err := service.NewEconomyService()
    if err != nil {
        panic(err)
    }
    defer cleanup()
    
    // Register commands
    commands.RegisterCommands(svc)
    
    // Server setup and start
    srv := server.DefaultConfig().New()
    srv.Listen()
    
    // Auto-register players on join
    for p := range srv.Accept() {
        svc.RegisterUser(context.Background(), p.UUID(), p.Name())
    }
}
```

### 3. Configuration

Database file (`foo.db`) will be created automatically. Player balance data is stored here.

## Features

- **Balance Management**: Check and set player balances
- **Transfer System**: Safe money transfers between players
- **Leaderboard**: Player rankings by balance
- **Auto Registration**: Automatic new player registration
- **Error Handling**: User-friendly error messages

## Requirements

- Go 1.21+
- df-mc/dragonfly framework
- SQLite (auto-setup)

## 日本語ドキュメント

日本語のドキュメントは[こちら](README_JP.md)をご覧ください。

## License

MIT License