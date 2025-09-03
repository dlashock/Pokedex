package main

import (
	"fmt"
	"os"
	"pokedexcli/internal/api"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func createCommandMap() map[string]cliCommand {
	commands := map[string]cliCommand{}
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    func() error { return commandHelp(commands) },
	}
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	commands["map"] = cliCommand{
		name:        "map",
		description: "Display a list of the location areas in the Pokemon games",
		callback:    commandMap,
	}
	return commands
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(commands map[string]cliCommand) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap() error {
	areas, err := api.ApiRequest()
	if err != nil {
		return fmt.Errorf("Error making API call: %w", err)
	}

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}
	return nil
}
