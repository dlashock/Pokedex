package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"pokedexcli/internal/api"
	"strings"
)

const DefaultSaveFileName = ".pokedex.json"

// StorageManager handles persistent storage of the user's Pokédex
type StorageManager struct {
	filePath string
}

// NewStorageManager creates a new storage manager with the specified file path
func NewStorageManager(filePath string) *StorageManager {
	return &StorageManager{
		filePath: filePath,
	}
}

// GetDefaultSavePath returns the default save path in the user's home directory
func GetDefaultSavePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, DefaultSaveFileName), nil
}

// PromptForSavePath prompts the user to choose where to save their Pokédex
func PromptForSavePath() (string, error) {
	defaultPath, err := GetDefaultSavePath()
	if err != nil {
		// Fallback to current directory if home dir fails
		defaultPath = DefaultSaveFileName
	}

	fmt.Printf("Where would you like to save your Pokédex? (Press Enter for default: %s)\n", defaultPath)
	fmt.Print("Save path: ")
	
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read user input: %w", err)
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultPath, nil
	}

	// Expand ~ to home directory if present
	if strings.HasPrefix(input, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			input = filepath.Join(homeDir, input[2:])
		}
	}

	// Ensure the directory exists
	dir := filepath.Dir(input)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	return input, nil
}

// LoadPokedex loads the Pokédex from the file system
func (sm *StorageManager) LoadPokedex() (map[string]api.Pokemon, error) {
	// Check if file exists
	if _, err := os.Stat(sm.filePath); os.IsNotExist(err) {
		// File doesn't exist, return empty Pokédex
		return make(map[string]api.Pokemon), nil
	}

	// Read the file
	data, err := os.ReadFile(sm.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Pokédex file: %w", err)
	}

	// Handle empty file
	if len(data) == 0 {
		return make(map[string]api.Pokemon), nil
	}

	// Parse JSON
	var pokedex map[string]api.Pokemon
	if err := json.Unmarshal(data, &pokedex); err != nil {
		// File is corrupted, create backup and start fresh
		backupPath := sm.filePath + ".backup"
		if backupErr := os.WriteFile(backupPath, data, 0644); backupErr == nil {
			fmt.Printf("Warning: Pokédex file was corrupted. Backed up to %s\n", backupPath)
		}
		fmt.Println("Starting with a fresh Pokédex...")
		return make(map[string]api.Pokemon), nil
	}

	if pokedex == nil {
		pokedex = make(map[string]api.Pokemon)
	}

	fmt.Printf("Loaded %d Pokémon from your saved Pokédex!\n", len(pokedex))
	return pokedex, nil
}

// SavePokedex saves the Pokédex to the file system using atomic write
func (sm *StorageManager) SavePokedex(pokedex map[string]api.Pokemon) error {
	// Marshal to JSON with pretty formatting
	data, err := json.MarshalIndent(pokedex, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Pokédex to JSON: %w", err)
	}

	// Create temporary file for atomic write
	tempPath := sm.filePath + ".tmp"
	
	// Write to temporary file
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, sm.filePath); err != nil {
		// Clean up temp file if rename fails
		os.Remove(tempPath)
		return fmt.Errorf("failed to save Pokédex file: %w", err)
	}

	return nil
}

// FileExists checks if the storage file exists
func (sm *StorageManager) FileExists() bool {
	_, err := os.Stat(sm.filePath)
	return !os.IsNotExist(err)
}

// GetFilePath returns the current file path
func (sm *StorageManager) GetFilePath() string {
	return sm.filePath
}

// ValidatePath ensures the directory for the file path exists
func ValidatePath(filePath string) error {
	dir := filepath.Dir(filePath)
	
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Try to create directory
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	
	// Test write permissions by creating a temporary file
	tempFile := filepath.Join(dir, ".pokedex_test")
	if err := os.WriteFile(tempFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("no write permission in directory %s: %w", dir, err)
	}
	
	// Clean up test file
	os.Remove(tempFile)
	return nil
}