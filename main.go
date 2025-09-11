package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	//Create map for all possible commands
	commands := createCommandMap()

	//Initialize input buffer
	rl, err := readline.New("Pokedex > ")
	if err != nil {
		fmt.Println("Error initializing readline:", err)
		os.Exit(1)
	}

	//Begin accepting input
	for {
		line, err := rl.Readline()
		if err != nil {
			// Handle EOF (Ctrl+D) gracefully
			fmt.Println("\nClosing the Pokedex... Goodbye!")
			break
		}

		//Clean provided input by separating by whitespace and forcing lowecase
		words := cleanInput(line)

		//Ask for input if none was provided
		if len(words) == 0 {
			continue
		}

		//Check if the command exists in the command map and execute if so
		command, exists := commands[words[0]]
		if exists {
			// Handle special case for inspect command with -sprite flag
			if words[0] == "inspect" && len(words) > 2 && words[len(words)-1] == "-sprite" {
				// Remove the -sprite flag and join the remaining words for Pokemon name
				pokemonName := strings.Join(words[1:len(words)-1], "-")
				err := commandInspectSprite(pokemonName, commands)
				if err != nil {
					fmt.Printf("An error has occurred: %s\n", err)
				}
			} else {
				// Join all words after the command with dashes for multi-word Pokemon names
				var arg string
				if len(words) > 1 {
					arg = strings.Join(words[1:], "-")
				}

				err := command.callback(arg)
				if err != nil {
					fmt.Printf("An error has occurred: %s\n", err)
				}
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
