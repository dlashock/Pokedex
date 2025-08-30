package main

import (
	"strings"
)

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	//cleanWords := []string
	//for _, word := range words {
	//}
	return words
}
