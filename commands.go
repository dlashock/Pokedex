package main

import (
	"fmt"
	"os"
	"pokedexcli/internal/api"
	"pokedexcli/internal/pokecache"
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
		description: "Displays a help message",
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
		description: "Display a list of the next 20 location areas in the Pokemon games",
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
		description: "Display a list of the previous 20 location areas in the Pokemon games",
		callback:    func(arg string) error { return commandExplore(arg, commands) },
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
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(_ string, commands map[string]cliCommand) error {
	areas, err := api.ApiRequest(commands["map"].config.next, commands["map"].config.pokecache)
	if err != nil {
		return fmt.Errorf("Error making API call: %w", err)
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
			fmt.Print("You are already on the first page.")
			return nil
		}
	}

	areas, err := api.ApiRequest(commands["mapb"].config.previous, commands["mapb"].config.pokecache)
	if err != nil {
		return fmt.Errorf("Error making API call: %w", err)
	}

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}

	commands["map"].config.next = areas.Next
	commands["mapb"].config.previous = areas.Previous

	return nil
}

func commandExplore(arg string, commands map[string]cliCommand) error {
	fmt.Printf("Exploring %s...\n", arg)
	return nil
}
