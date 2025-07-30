# DF Economy

Economy system library for Minecraft Bedrock Edition (**df-mc/dragonfly** framework)

## Overview

DF Economy provides a flexible currency system for Minecraft servers with multi-database support, player-to-player transactions, balance checking, and leaderboard features. The system supports SQLite, MySQL, and PostgreSQL databases.

## In-Game Commands

| Command | Description | Example |
| --- | --- | --- |
| `/economy` | Show command help | `/economy` |
| `/economy balance [player]` | Display balance | `/economy balance` or `/economy balance Steve` |
| `/economy pay <player> <amount>` | Send money to another player | `/economy pay Steve 100` |
| `/economy set <player> <amount>` | Set player balance (configurable) | `/economy set Steve 1000` |
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
    "context"
    "github.com/df-mc/dragonfly/server"
    "github.com/skuralll/df-permission/permission"
    "github.com/skuralll/dfeconomy/dragonfly/commands"
    "github.com/skuralll/dfeconomy/economy/config"
    "github.com/skuralll/dfeconomy/economy/service"
)

func main() {
    // Create economy service with configuration
    cfg := config.Config{
        DBType:         "sqlite",          // or "mysql", "postgres"
        DBDSN:          "./economy.db",    // database connection string
        DefaultBalance: 100.0,             // starting balance for new players
        EnableSetCmd:   false,             // enable /economy set command (default: false)
    }
    
    // Initialize permission manager
    // Pass nil if permission management is not needed
    pMgr := permission.NewManager()
    
    svc, cleanup, err := service.NewEconomyService(cfg, pMgr)
    if err != nil {
        panic(err)
    }
    defer cleanup()
    
    // Register commands
    commands.RegisterCommands(svc, cfg)
    
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

The economy system supports multiple database backends:

#### SQLite (Default)
```go
cfg := config.Config{
    DBType: "sqlite",
    DBDSN:  "./economy.db",
    DefaultBalance: 100.0,
    EnableSetCmd: false,
}
```

#### MySQL
```go
cfg := config.Config{
    DBType: "mysql",
    DBDSN:  "user:password@tcp(localhost:3306)/economy?charset=utf8mb4&parseTime=True&loc=Local",
    DefaultBalance: 100.0,
    EnableSetCmd: false,
}
```

#### PostgreSQL
```go
cfg := config.Config{
    DBType: "postgres",
    DBDSN:  "host=localhost user=user password=password dbname=economy port=5432 sslmode=disable",
    DefaultBalance: 100.0,
    EnableSetCmd: false,
}
```

Database tables and schemas are automatically created on startup.

#### Enable Set Command (Optional)
To enable the `/economy set` command for balance management:
```go
cfg := config.Config{
    DBType: "sqlite",
    DBDSN:  "./economy.db",
    DefaultBalance: 100.0,
    EnableSetCmd: true,  // Enable set command
}
```

**Note**: The set command is disabled by default for security reasons.

## Features

- **Multi-Database Support**: SQLite, MySQL, and PostgreSQL support
- **Balance Management**: Check and set player balances
- **Transfer System**: Safe money transfers between players
- **Leaderboard**: Player rankings by balance
- **Auto Registration**: Automatic new player registration with configurable starting balance
- **Error Handling**: User-friendly error messages with proper validation
- **CGO-Free**: Pure Go implementation for all database drivers
- **Transaction Safety**: ACID compliance with proper rollback handling
- **Command Control**: Configurable command availability for enhanced security

## Requirements

- Go 1.24+
- df-mc/dragonfly framework
- Database: SQLite (auto-setup), MySQL, or PostgreSQL

## 日本語ドキュメント

日本語のドキュメントは[こちら](README_JP.md)をご覧ください。

## License

MIT License