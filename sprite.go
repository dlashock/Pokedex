package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"time"

	"pokedexcli/internal/api"

	"github.com/qeesung/image2ascii/convert"
)

// displaySprite fetches and displays the ASCII art for a given Pokémon
func displaySprite(pokemon api.Pokemon) error {
	// Get the front default sprite URL
	spriteURL := pokemon.Sprites.FrontDefault
	if spriteURL == "" {
		return fmt.Errorf("no sprite available for %s", pokemon.Name)
	}

	// Download the sprite image
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(spriteURL)
	if err != nil {
		return fmt.Errorf("failed to download sprite: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download sprite: HTTP %d", resp.StatusCode)
	}

	// Read the image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read sprite data: %w", err)
	}

	// Decode the image
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return fmt.Errorf("failed to decode sprite image: %w", err)
	}

	// Configure the ASCII converter
	convertOptions := convert.DefaultOptions
	convertOptions.FixedWidth = 40
	convertOptions.FixedHeight = 20
	convertOptions.Colored = true

	// Create converter
	converter := convert.NewImageConverter()

	// Convert image to ASCII
	asciiArt := converter.Image2ASCIIString(img, &convertOptions)

	// Display the ASCII art
	fmt.Println(asciiArt)

	return nil
}

// displaySpriteByName displays a sprite for a Pokémon by name from the user's Pokédex
func displaySpriteByName(pokemonName string, pokedex map[string]api.Pokemon) error {
	pokemon, exists := pokedex[pokemonName]
	if !exists {
		return fmt.Errorf("you have not caught %s yet", pokemonName)
	}

	return displaySprite(pokemon)
}
