# DF Economy Project

Economic system for Minecraft Dragonfly server

## About Dragonfly
- Asynchronous Minecraft: Bedrock Edition server written in Go
- Extensibility-focused design, used as a library
- Struct-based command system

## Command System Basics
- Implement `cmd.Runnable` interface
- Struct fields become command parameters
- Use `cmd.SubCommand` to define subcommands
- Use `cmd.Optional[T]` for optional arguments

## Project Structure
```
dragonfly/commands/
├── base.go        # BaseCommand definition
├── register.go    # Command registration
├── economy.go     # Main command
├── balance.go     # Balance display
├── set.go         # Balance setting
├── top.go         # Rankings
└── pay.go         # Money transfer
```

## Important Design Decisions
- Commands are split by file
- BaseCommand provides shared service reference
- Error handling is implemented individually per command

## Information Retrieval
When detailed Dragonfly information is needed, use the DeepwikiMCP server:

- **Overview**: Use `mcp__deepwiki__ask_question` to get basic info about df-mc/dragonfly
- **Structure**: Use `mcp__deepwiki__read_wiki_structure` to check documentation structure
- **Command System**: Ask "How does the command system work?" for detailed specifications
- **Entity System**: Ask "How do entities work?" for implementation details
- **World Management**: Ask "How does world management work?" for world-related specifications