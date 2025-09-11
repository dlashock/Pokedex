package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
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
	pokedex   map[string]api.Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(arg string) error
	config      *config
}

func createCommandMap() map[string]cliCommand {
	freshCache := pokecache.NewCache(2 * time.Minute)
	commands := map[string]cliCommand{}

	sharedConfig := &config{
		next:      "https://pokeapi.co/api/v2/location-area/",
		previous:  "",
		pokecache: freshCache,
		pokedex:   make(map[string]api.Pokemon),
	}

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
		description: "See details of a Pokemon you have caught. Takes the name of a Pokemon as an argument",
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
	fmt.Printf("Throwing a Pokeball at %s...\n", arg)

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

	catch := catchAttempt(pokemon.BaseExperience)

	if catch {
		commands["catch"].config.pokedex[arg] = pokemon
		fmt.Printf("%s was caught!\n", arg)
		fmt.Println("You may now inspect it with the inspect command")
	} else {
		fmt.Printf("%s escaped!\n", arg)
	}

	return nil
}

func catchAttempt(baseEXP int) bool {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	catch_rate := 1.95 - 0.279*math.Log(float64(baseEXP))
	rand_num := rng.Float64()
	return catch_rate > rand_num
}

func commandInspect(arg string, commands map[string]cliCommand) error {
	val, exists := commands["inspect"].config.pokedex[arg]
	if !exists {
		fmt.Printf("You have not caught %s yet!\n", arg)
		return nil
	}
	printPokemon(val)
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
	fmt.Println("Your Pokedex:")
	for _, pokemon := range commands["pokedex"].config.pokedex {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
}
