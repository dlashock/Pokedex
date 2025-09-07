# Pokédex CLI

A command-line Pokédex application built with Go that provides an interactive interface for exploring Pokémon location data through the PokéAPI.

## Features

- **Interactive REPL**: Command-line interface with prompt-based navigation
- **Location Mapping**: Browse through Pokémon location areas with pagination
- **Intelligent Caching**: Built-in caching system with automatic TTL cleanup for fast response times
- **Location Exploration**: Explore specific areas to discover Pokémon encounters
- **Thread-Safe**: Concurrent access support with mutex-protected operations

## Commands

- `help` - Display available commands and their descriptions
- `exit` - Exit the Pokédex application
- `map` - Show the next 20 location areas
- `mapb` - Show the previous 20 location areas
- `explore <area-name>` - Explore a specific location area to find Pokémon

## Installation & Usage

### Prerequisites
- Go 1.25.0 or higher

### Build and Run
```bash
# Clone the repository
git clone <repository-url>
cd Pokedex

# Build the application
go build -o pokedexcli

# Run the application
./pokedexcli
```

### Example Usage
```
Pokedex > help
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
map: Display a list of the next 20 location areas in the Pokemon games
mapb: Display a list of the previous 20 location areas in the Pokemon games
explore: Explore a specific location area

Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
...

Pokedex > explore pastoria-city-area
Exploring pastoria-city-area...
Found Pokemon:
 - tentacool
 - tentacruel
 - magikarp
 ...

Pokedex > exit
Closing the Pokedex... Goodbye!
```

## Architecture

### Core Components
- **main.go** - Entry point with REPL loop for command processing
- **commands.go** - Command system using map-based dispatcher with state management
- **input.go** - Input processing and normalization utilities
- **internal/api/** - HTTP client for PokéAPI integration
- **internal/pokecache/** - Thread-safe caching system with TTL

### Caching System
The application features a sophisticated caching layer that:
- Stores API responses in memory with timestamps
- Automatically removes expired entries (5-minute TTL)
- Runs background cleanup processes via goroutines
- Provides thread-safe concurrent access with mutex protection
- Significantly reduces API calls and improves response times

### Command System
- Commands accept parameters through a flexible callback system
- State management enables bi-directional pagination through location areas
- Shared cache instance across all commands for optimal performance
- Extensible architecture for adding new commands

## Development

### Testing
```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run specific package tests
go test ./internal/pokecache/
go test ./internal/api/
```

### Module Management
```bash
# Tidy dependencies
go mod tidy

# Download dependencies
go mod download
```

## API Integration

This application integrates with the [PokéAPI](https://pokeapi.co/) to fetch:
- Location area lists with pagination
- Detailed location area data with Pokémon encounters
- Comprehensive error handling for network and parsing issues

## Project Status

This is an active learning project following the Boot.dev curriculum. Current functionality includes location browsing and exploration, with caching optimization for performance.

## License

This project is for educational purposes as part of the Boot.dev Go course.