package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"pokedexcli/internal/api"
	"pokedexcli/internal/pokecache"
	"sort"
	"time"
)

// Global random generator for consistent seeding
var globalRng *rand.Rand

func init() {
	globalRng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type config struct {
	next           string
	previous       string
	pokecache      *pokecache.Cache
	pokedex        map[string]api.Pokemon
	storageManager *StorageManager
}

type cliCommand struct {
	name        string
	description string
	callback    func(arg string) error
	config      *config
}

func createCommandMap() map[string]cliCommand {
	commands := map[string]cliCommand{}
	sharedConfig := createSharedConfig()

	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays all available commands and information about what they do",
		callback:    func(arg string) error { return commandHelp(arg, commands) },
		config:      sharedConfig,
	}
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
		config:      sharedConfig,
	}
	commands["map"] = cliCommand{
		name:        "map",
		description: "Display a list of the next 20 location areas in the Pokemon games.",
		callback:    func(arg string) error { return commandMap(arg, commands) },
		config:      sharedConfig,
	}
	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Display a list of the previous 20 location areas in the Pokemon games",
		callback:    func(arg string) error { return commandMapb(arg, commands) },
		config:      sharedConfig,
	}
	commands["explore"] = cliCommand{
		name:        "explore",
		description: "Display a list of Pokemon in the provided area. Accepts a single location area as an argument",
		callback:    func(arg string) error { return commandExplore(arg, commands) },
		config:      sharedConfig,
	}
	commands["catch"] = cliCommand{
		name:        "catch",
		description: "Try to catch a Pokemon! Takes a Pokemon name as an an argument",
		callback:    func(arg string) error { return commandCatch(arg, commands) },
		config:      sharedConfig,
	}
	commands["inspect"] = cliCommand{
		name:        "inspect",
		description: "See details of a Pokemon you have caught. Takes the name of a Pokemon as an argument. Add -sprite flag to also display ASCII sprite",
		callback:    func(arg string) error { return commandInspect(arg, commands) },
		config:      sharedConfig,
	}
	commands["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "See the list of Pokemon you have caught",
		callback: func(arg string) error {
			commandPokedex(commands)
			return nil
		},
		config: sharedConfig,
	}
	commands["sprite"] = cliCommand{
		name:        "sprite",
		description: "Display the ASCII sprite of a Pokemon you have caught",
		callback:    func(arg string) error { return commandSprite(arg, commands) },
		config:      sharedConfig,
	}
	return commands
}

func commandExit(_ string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ string, commands map[string]cliCommand) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	keys := make([]string, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}

	// Sort the keys
	sort.Strings(keys)

	// Iterate over the sorted keys and print values
	for _, k := range keys {
		fmt.Printf("%s: %s\n", commands[k].name, commands[k].description)
	}
	return nil
}

func commandMap(_ string, commands map[string]cliCommand) error {
	body, err := api.ApiRequest(commands["map"].config.next, commands["map"].config.pokecache)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %w", err)
	}

	var areas api.LocationArea
	if err := unmarshalJSON(body, &areas); err != nil {
		return err
	}

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}

	commands["map"].config.next = areas.Next
	commands["mapb"].config.previous = areas.Previous

	return nil
}

func commandMapb(_ string, commands map[string]cliCommand) error {
	if commands["mapb"].config.previous == "" {
		fmt.Println("You are already on the first page.")
		return nil
	}

	body, err := api.ApiRequest(commands["mapb"].config.previous, commands["mapb"].config.pokecache)
	if err != nil {
		return fmt.Errorf("failed to fetch previous location areas: %w", err)
	}

	var areas api.LocationArea
	if err := unmarshalJSON(body, &areas); err != nil {
		return err
	}

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}

	commands["map"].config.next = areas.Next
	commands["mapb"].config.previous = areas.Previous

	return nil
}

func commandExplore(arg string, commands map[string]cliCommand) error {
	if err := validateRequiredArg(arg, "explore"); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	fmt.Printf("Exploring %s...\n", arg)

	body, err := api.ApiRequest(BaseLocationAreaURL+arg, commands["explore"].config.pokecache)
	if err != nil {
		if apiErr, ok := err.(api.APIError); ok && apiErr.IsNotFound() {
			return fmt.Errorf("area '%s' does not exist. Please check spelling and try again", arg)
		}
		return fmt.Errorf("failed to explore area: %w", err)
	}

	var area api.Area
	if err := unmarshalJSON(body, &area); err != nil {
		return err
	}

	if len(area.PokemonEncounters) > 0 {
		fmt.Println("Found Pokemon:")
		for _, pokemon := range area.PokemonEncounters {
			fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
		}
	} else {
		fmt.Println("No Pokemon found in this area.")
	}
	return nil
}

func commandCatch(arg string, commands map[string]cliCommand) error {
	if err := validateRequiredArg(arg, "catch"); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", arg)

	body, err := api.ApiRequest(BasePokemonURL+arg, commands["catch"].config.pokecache)
	if err != nil {
		if apiErr, ok := err.(api.APIError); ok && apiErr.IsNotFound() {
			return fmt.Errorf("pokemon '%s' does not exist. Please check spelling and try again", arg)
		}
		return fmt.Errorf("failed to fetch pokemon data: %w", err)
	}

	var pokemon api.Pokemon
	if err := unmarshalJSON(body, &pokemon); err != nil {
		return err
	}

	if catchAttempt(pokemon.BaseExperience) {
		commands["catch"].config.pokedex[arg] = pokemon
		fmt.Printf("%s was caught!\n", arg)
		fmt.Println("You may now inspect it with the inspect command")
		
		// Save the updated Pokédex
		if err := commands["catch"].config.storageManager.SavePokedex(commands["catch"].config.pokedex); err != nil {
			fmt.Printf("Warning: Failed to save Pokédex: %v\n", err)
			fmt.Println("Your progress is still safe in memory for this session.")
		}
	} else {
		fmt.Printf("%s escaped!\n", arg)
	}

	return nil
}

func catchAttempt(baseExp int) bool {
	catchRate := CatchRateConstantA - CatchRateConstantB*math.Log(float64(baseExp))
	randomValue := globalRng.Float64()
	return catchRate > randomValue
}

func commandInspect(arg string, commands map[string]cliCommand) error {
	if err := validateRequiredArg(arg, "inspect"); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	pokemon, exists := commands["inspect"].config.pokedex[arg]
	if !exists {
		fmt.Printf("You have not caught %s yet!\n", arg)
		return nil
	}

	printPokemon(pokemon)
	return nil
}

func printPokemon(mon api.Pokemon) {
	fmt.Printf("Name: %s\n", mon.Name)
	fmt.Printf("Height: %d\n", mon.Height)
	fmt.Printf("Weight: %d\n", mon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range mon.Stats {
		fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, t := range mon.Types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}
}

func commandPokedex(commands map[string]cliCommand) {
	pokedex := commands["pokedex"].config.pokedex
	if len(pokedex) == 0 {
		fmt.Println("Your Pokedex is empty. Try catching some Pokemon first!")
		return
	}

	fmt.Println("Your Pokedex:")
	for _, pokemon := range pokedex {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
}

func commandSprite(arg string, commands map[string]cliCommand) error {
	if err := validateRequiredArg(arg, "sprite"); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	pokedex := commands["sprite"].config.pokedex
	return displaySpriteByName(arg, pokedex)
}

func commandInspectSprite(arg string, commands map[string]cliCommand) error {
	if err := validateRequiredArg(arg, "inspect"); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	pokemon, exists := commands["inspect"].config.pokedex[arg]
	if !exists {
		fmt.Printf("You have not caught %s yet!\n", arg)
		return nil
	}

	// Print normal Pokemon info first
	printPokemon(pokemon)

	return displaySprite(pokemon)
}
