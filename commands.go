package main

import (
	"encoding/json"
	"fmt"
	"os"
	"pokedexcli/internal/api"
	"pokedexcli/internal/pokecache"
	"sort"
	"strings"
	"time"
)

type config struct {
	next      string
	previous  string
	pokecache *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(arg string) error
	config      *config
}

func createCommandMap() map[string]cliCommand {
	freshCache := pokecache.NewCache(5 * time.Minute)
	commands := map[string]cliCommand{}
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays all available commands and information about what they do",
		callback:    func(arg string) error { return commandHelp(arg, commands) },
		config:      &config{pokecache: freshCache},
	}
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
		config:      &config{pokecache: freshCache},
	}
	commands["map"] = cliCommand{
		name:        "map",
		description: "Display a list of the next 20 location areas in the Pokemon games.",
		callback:    func(arg string) error { return commandMap(arg, commands) },
		config: &config{
			next:      "https://pokeapi.co/api/v2/location-area/",
			pokecache: freshCache,
		},
	}
	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Display a list of the previous 20 location areas in the Pokemon games",
		callback:    func(arg string) error { return commandMapb(arg, commands) },
		config: &config{
			previous:  "",
			pokecache: freshCache,
		},
	}
	commands["explore"] = cliCommand{
		name:        "explore",
		description: "Display a list of Pokemon in the provided area. Accepts a single location area as an argument",
		callback:    func(arg string) error { return commandExplore(arg, commands) },
		config: &config{
			pokecache: freshCache,
		},
	}
	commands["catch"] = cliCommand{
		name:        "catch",
		description: "Try to catch a Pokemon in the current area",
		callback:    func(arg string) error { return commandCatch(arg, commands) },
		config: &config{
			pokecache: freshCache,
		},
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
		return fmt.Errorf("Error making API call: %w", err)
	}

	var areas api.LocationArea
	if err := json.Unmarshal(body, &areas); err != nil {
		return fmt.Errorf("Error unmarshalling JSON: %v", err)
	}

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}

	commands["map"].config.next = areas.Next
	commands["mapb"].config.previous = areas.Previous

	return nil
}

func commandMapb(_ string, commands map[string]cliCommand) error {
	if commands["mapb"].config != nil {
		if commands["mapb"].config.previous == "" {
			fmt.Println("You are already on the first page.")
			return nil
		}
	}
	var areas api.LocationArea

	body, err := api.ApiRequest(commands["mapb"].config.previous, commands["mapb"].config.pokecache)
	if err != nil {
		return fmt.Errorf("Error making API call: %w", err)
	}

	if err := json.Unmarshal(body, &areas); err != nil {
		return fmt.Errorf("Error unmarshalling JSON: %w", err)
	}

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}

	commands["map"].config.next = areas.Next
	commands["mapb"].config.previous = areas.Previous

	return nil
}

func commandExplore(arg string, commands map[string]cliCommand) error {
	if strings.TrimSpace(arg) == "" {
		fmt.Print("Please provide a location to check for Pokemon")
		return nil
	}
	fmt.Printf("Exploring %s...\n", arg)

	var area api.Area
	body, err := api.ApiRequest("https://pokeapi.co/api/v2/location-area/"+arg, commands["explore"].config.pokecache)
	if err != nil {
		if strings.Contains(err.Error(), "status code: 404") {
			return fmt.Errorf("Area '%s' does not exist. Please check spelling and try again", arg)
		}
		return fmt.Errorf("Error exploring area: %w", err)
	}

	if err := json.Unmarshal(body, &area); err != nil {
		return fmt.Errorf("Error unmarshalling JSON: %w", err)
	}

	if area.PokemonEncounters != nil {
		fmt.Print("Found Pokemon:\n")
		for _, pokemon := range area.PokemonEncounters {
			fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
		}
	}
	return nil
}

func commandCatch(arg string, commands map[string]cliCommand) error {
	fmt.Printf("Throwing a PokeBall at %s...\n", arg)

	var pokemon api.Pokemon
	body, err := api.ApiRequest("https://pokeapi.co/api/v2/pokemon/"+arg, commands["catch"].config.pokecache)
	if err != nil {
		if strings.Contains(err.Error(), "status code: 404") {
			return fmt.Errorf("Pokemon '%s' does not exist. Please check spelling and try again", arg)
		}
		return fmt.Errorf("Error trying to catch %s: %w", arg, err)
	}

	if err := json.Unmarshal(body, &pokemon); err != nil {
		return fmt.Errorf("Error unmarshalling JSON: %w", err)
	}

	catch := catchChance(pokemon.BaseExperience)

	if catch {
		fmt.Printf("%s was caught!\n", arg)
	} else {
		fmt.Printf("%s escaped!\n", arg)
	}

	return nil
}

func catchChance(baseEXP int) bool {
	return false
}
