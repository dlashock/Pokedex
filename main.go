package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
	words := cleanInput("My Name is pOOp butt . ")
	for _, word := range words {
		fmt.Println(word)
	}
	fmt.Println(len(words))
}
