# dfeconomy

A lightweight economy library for the **df-mc/dragonfly** Bedrock server framework.


## Overview
dfeconomy adds a simple in-game currency system that can be queried, modified and extended through both chat commands and a concise Go API.  


## Features
| Feature               | Details                                                        |
| --------------------- | -------------------------------------------------------------- |
| **Minecraft Support** | Tested on **Bedrock 1.21.81**                                  |
| **Storage**           | Built-in **SQLite** backend (in-memory backend also available) |

## Commands
| Command                        | Description                         |
| ------------------------------ | ----------------------------------- |
| `/money`                       | Show command help                   |
| `/money help`                  | Show command help                   |
| `/money set <player> <amount>` | Set the target player’s balance     |
| `/money pay <player> <amount>` | Transfer money to the target player |
| `/money top <page>`            | Show balance leaderboard            |

## API
The core package exposes a minimal interface for server-side plugins:

```go
// Retrieve a player’s balance
economy.Balance(srv, p)

// Add to a player’s balance
economy.Deposit(srv, p, amount)

// Subtract from a player’s balance
economy.Withdraw(srv, p, amount)

// Set an exact balance
economy.Set(srv, p, amount)

// Fetch a paginated leaderboard
economy.Top(srv, page, pageSize) ([]economy.Entry, error)