package main

import (
	"encoding/json"
	"fmt"
	"pokedexcli/internal/api"
	"pokedexcli/internal/pokecache"
	"strings"
)

// validateRequiredArg validates that a required argument is provided and not empty
func validateRequiredArg(arg, commandName string) error {
	if strings.TrimSpace(arg) == "" {
		return NewValidationError("argument", fmt.Sprintf("%s command requires an argument", commandName))
	}
	return nil
}

// unmarshalJSON is a generic helper for JSON unmarshaling with better error handling
func unmarshalJSON[T any](data []byte, target *T) error {
	if err := json.Unmarshal(data, target); err != nil {
		return NewJSONError(err.Error(), string(data))
	}
	return nil
}

// createSharedConfig creates the shared configuration for all commands
func createSharedConfig() *config {
	// Initialize storage manager
	storageManager, pokedex := initializeStorage()
	
	return &config{
		next:           BaseLocationAreaURL,
		previous:       "",
		pokecache:      newCache(),
		pokedex:        pokedex,
		storageManager: storageManager,
	}
}

// Helper function to create cache with consistent configuration
func newCache() *pokecache.Cache {
	return pokecache.NewCache(CacheTTL)
}

// initializeStorage sets up the storage manager and loads existing Pokédex data
func initializeStorage() (*StorageManager, map[string]api.Pokemon) {
	var savePath string
	var err error
	
	// Try to get default path first
	defaultPath, err := GetDefaultSavePath()
	if err != nil {
		defaultPath = DefaultSaveFileName
	}
	
	// Check if default file exists
	defaultManager := NewStorageManager(defaultPath)
	if defaultManager.FileExists() {
		savePath = defaultPath
	} else {
		// No existing file, prompt user for save location
		fmt.Println("Welcome to the Pokédex! Let's set up your save file.")
		savePath, err = PromptForSavePath()
		if err != nil {
			fmt.Printf("Error setting up save path: %v\n", err)
			fmt.Println("Falling back to current directory...")
			savePath = DefaultSaveFileName
		}
		
		// Validate the chosen path
		if err := ValidatePath(savePath); err != nil {
			fmt.Printf("Warning: Cannot write to %s: %v\n", savePath, err)
			fmt.Println("Using current directory instead...")
			savePath = DefaultSaveFileName
		}
	}
	
	// Create storage manager
	storageManager := NewStorageManager(savePath)
	
	// Load existing Pokédex data
	pokedex, err := storageManager.LoadPokedex()
	if err != nil {
		fmt.Printf("Error loading Pokédex: %v\n", err)
		fmt.Println("Starting with empty Pokédex...")
		pokedex = make(map[string]api.Pokemon)
	}
	
	fmt.Printf("Pokédex will be saved to: %s\n", savePath)
	return storageManager, pokedex
}
