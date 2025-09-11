# Pokédex CLI

A feature-rich command-line Pokédex game built with Go that provides an interactive interface for exploring Pokémon locations, catching Pokémon, and building your collection through the PokéAPI.

## Features

- **Interactive REPL**: Command-line interface with prompt-based navigation
- **Location Mapping**: Browse through Pokémon location areas with pagination
- **Pokémon Exploration**: Discover which Pokémon inhabit specific areas
- **Pokémon Catching**: Catch Pokémon with probability-based mechanics
- **Collection Management**: Build and view your personal Pokédex
- **Intelligent Caching**: Built-in caching system with automatic TTL cleanup for blazing-fast response times
- **Thread-Safe**: Concurrent access support with mutex-protected operations

## Commands

- `help` - Display available commands and their descriptions
- `exit` - Exit the Pokédex application
- `map` - Show the next 20 location areas
- `mapb` - Show the previous 20 location areas
- `explore <area-name>` - Explore a specific location area to find Pokémon
- `catch <pokemon-name>` - Attempt to catch a Pokémon (probability-based)
- `inspect <pokemon-name>` - View detailed stats of a caught Pokémon
- `pokedex` - Display all Pokémon you've caught

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

catch: Try to catch a Pokemon! Takes a Pokemon name as an an argument
exit: Exit the Pokedex
explore: Display a list of Pokemon in the provided area. Accepts a single location area as an argument
help: Displays all available commands and information about what they do
inspect: See details of a Pokemon you have caught. Takes the name of a Pokemon as an argument
map: Display a list of the next 20 location areas in the Pokemon games.
mapb: Display a list of the previous 20 location areas in the Pokemon games
pokedex: See the list of Pokemon you have caught

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
 - gyarados

Pokedex > catch magikarp
Throwing a Pokeball at magikarp...
magikarp was caught!
You may now inspect it with the inspect command

Pokedex > inspect magikarp
Name: magikarp
Height: 9
Weight: 100
Stats:
 -hp: 20
 -attack: 10
 -defense: 55
 -special-attack: 15
 -special-defense: 20
 -speed: 80
Types:
 - water

Pokedex > pokedex
Your Pokedex:
 - magikarp

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

### Game Mechanics
The application features sophisticated Pokémon game mechanics:
- **Probabilistic Catching**: Uses logarithmic formula based on Pokémon base experience (95% chance for weakest, 15% for strongest)
- **Persistent Collection**: Caught Pokémon are stored in your personal Pokédex throughout the session
- **Detailed Inspection**: View complete Pokémon stats including HP, attack, defense, types, height, and weight
- **Collection Management**: Track all caught Pokémon with the dedicated pokedex command

### Caching System
The application features a sophisticated caching layer that:
- Stores API responses in memory with timestamps
- Automatically removes expired entries (2-minute TTL)
- Runs background cleanup processes via goroutines
- Provides thread-safe concurrent access with mutex protection
- Significantly reduces API calls and improves response times

### Command System
- Commands accept parameters through a flexible callback system
- State management enables bi-directional pagination and persistent Pokémon collection
- Shared state across all commands for optimal performance and data consistency
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
- Individual Pokémon data including stats, types, and base experience
- Comprehensive error handling for network and parsing issues with user-friendly 404 messages

## Project Status

This is a completed learning project following the Boot.dev curriculum. Features a full Pokédex game experience including location exploration, Pokémon catching with realistic probability mechanics, collection management, and performance optimization through intelligent caching.

## License

This project is for educational purposes as part of the Boot.dev Go course.