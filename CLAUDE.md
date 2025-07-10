# DF Economy Project

Economic system for Minecraft Dragonfly server

## About Dragonfly
- Asynchronous Minecraft: Bedrock Edition server written in Go
- Extensibility-focused design, used as a library
- Struct-based command system with reflection-based argument parsing
- Client-side command integration (auto-completion, validation)

## Command System Basics
- Implement `cmd.Runnable` interface
- Struct fields become command parameters
- Use `cmd.SubCommand` to define subcommands
- Use `cmd.Optional[T]` for optional arguments
- Commands registered via `cmd.Register()` with `cmd.New()`
- Error handling through `cmd.Output` and direct player messaging

## Project Structure
```
df-economy/
├── dragonfly/commands/     # Command implementations
│   ├── base.go            # BaseCommand with shared utilities
│   ├── register.go        # Command registration
│   ├── economy.go         # Main help command
│   ├── balance.go         # Balance display (supports optional target)
│   ├── set.go             # Balance setting (admin command)
│   ├── top.go             # Rankings with pagination
│   └── pay.go             # Money transfer between players
├── economy/service/       # Business logic layer
│   └── service.go         # EconomyService with validation
├── internal/db/           # Database abstraction
│   ├── db.go              # DB interface
│   └── sqlite.go          # SQLite implementation
├── models/                # Data models
│   └── economy.go         # EconomyEntry struct
├── errors/                # Custom error types
│   └── errors.go          # Economy-specific errors
├── cmd/demo/              # Demo server implementation
│   └── main.go            # Server setup with economy integration
└── config.toml            # Dragonfly server configuration
```

## Architecture Design
- **Service Layer**: `EconomyService` handles business logic and validation
- **Database Layer**: Interface-based design with SQLite implementation  
- **Command Layer**: Dragonfly commands with async execution
- **Error Handling**: Custom error types with user-friendly messages
- **Async Processing**: Non-blocking command execution with timeout handling

## Command Implementation Patterns
- **BaseCommand**: Shared utilities for player validation, async execution, UUID resolution
- **Async Execution**: All database operations run asynchronously with 5-second timeout
- **Error Messaging**: Consistent error handling with colored chat messages
- **Immediate Feedback**: Commands provide instant feedback before async processing

## Database Operations
- **Balance Management**: Get, Set, Transfer operations with atomic transactions
- **User Registration**: Automatic user registration on first join
- **Leaderboards**: Paginated top balance queries
- **UUID Resolution**: Name-to-UUID mapping for player operations

## Security & Validation
- **Input Validation**: Amount validation (positive, minimum values)
- **Self-Transfer Prevention**: Cannot pay yourself
- **Timeout Handling**: All operations have 5-second timeout
- **Error Logging**: Internal errors logged with context

## Information Retrieval
When detailed Dragonfly information is needed, use the DeepwikiMCP server:

- **Overview**: Use `mcp__deepwiki__ask_question` to get basic info about df-mc/dragonfly
- **Structure**: Use `mcp__deepwiki__read_wiki_structure` to check documentation structure
- **Command System**: Ask "How does the command system work?" for detailed specifications
- **Entity System**: Ask "How do entities work?" for implementation details
- **World Management**: Ask "How does world management work?" for world-related specifications

## Development Notes
- Commands use struct-based parameter definition (Dragonfly pattern)
- All database operations are context-aware with timeout handling
- Error messages use Minecraft color codes (§a for success, §c for error)
- Service layer provides clean separation between commands and database
- Demo server in `cmd/demo/main.go` shows complete integration example

## Memories
- to memorize