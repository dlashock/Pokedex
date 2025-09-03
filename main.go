package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	//Create map for all possible commands
	commands := createCommandMap()

	//Initialize input buffer
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex > ")

	//Begin accepting input
	for scanner.Scan() {
		line := scanner.Text()

		//Clean provided input by separating by whitespace and forcing lowecase
		words := cleanInput(line)

		//Ask for input if none was provided
		if len(words) == 0 {
			fmt.Println("Please enter a command")
			fmt.Print("Pokedex > ")
			continue
		}

		//Check if the command exists in the command map and execute if so
		command, exists := commands[words[0]]
		if exists {
			err := command.callback()
			if err != nil {
				fmt.Printf("An error has occurred: %s\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}

		//Print line to start the loop over
		fmt.Print("\nPokedex > ")
	}
}
