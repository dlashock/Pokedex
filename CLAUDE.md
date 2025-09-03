# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Pokédex CLI application written in Go that provides an interactive command-line interface for exploring Pokémon data. The application uses a REPL (Read-Eval-Print Loop) pattern to accept user commands and interact with the PokéAPI.

## Commands

### Build and Run
```bash
go build -o pokedexcli
./pokedexcli
```

### Testing
```bash
go test
```

### Module Management
```bash
go mod tidy
go mod download
```

## Architecture

### Core Components

- **main.go**: Entry point with REPL loop that handles user input and command execution
- **commands.go**: Command system using a map-based dispatcher pattern with `cliCommand` struct
- **input.go**: Input processing utilities for cleaning and parsing user commands
- **internal/apicalls.go**: HTTP client for PokéAPI with structured response types

### Command System

The application uses a command pattern where:
- Commands are defined as `cliCommand` structs with name, description, and callback function
- `createCommandMap()` in commands.go registers all available commands
- The main loop dispatches commands by looking them up in the command map
- Currently implements: `help`, `exit`, and `map` commands

### API Integration

- Uses Go's standard `net/http` client to interact with PokéAPI (https://pokeapi.co/api/v2/)
- Defines structured types for API responses (e.g., `locationArea` struct)
- Error handling includes HTTP status code validation and JSON unmarshaling errors

### Input Processing

- `cleanInput()` function normalizes user input by converting to lowercase and splitting on whitespace
- Handles empty input gracefully by prompting user to enter a command

### Testing

- Uses Go's standard testing package
- Test file: `repl_test.go` contains unit tests for input processing functions
- Note: Test function name has typo (`TestCleanInpu` instead of `TestCleanInput`)

## File Structure

```
├── main.go              # Main entry point and REPL
├── commands.go          # Command definitions and handlers
├── input.go            # Input processing utilities
├── internal/
│   └── apicalls.go     # API client and response types
├── repl_test.go        # Unit tests
├── go.mod              # Go module definition
└── pokedexcli          # Compiled binary
```