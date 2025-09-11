package main

import "time"

// API Configuration
const (
	// PokéAPI base URLs
	BaseLocationAreaURL = "https://pokeapi.co/api/v2/location-area/"
	BasePokemonURL      = "https://pokeapi.co/api/v2/pokemon/"

	// Cache configuration
	CacheTTL = 2 * time.Minute

	// Pokémon catching probability constants
	// Based on logarithmic formula: catch_rate = CatchRateConstantA - CatchRateConstantB * ln(base_experience)
	// Calibrated for 95% catch rate at 36 base exp, 15% at 635 base exp
	CatchRateConstantA = 1.95
	CatchRateConstantB = 0.279

	// Base experience bounds for reference
	MinBaseExperience = 36  // Lowest base experience in PokéAPI
	MaxBaseExperience = 635 // Highest base experience in PokéAPI

	// Expected catch rates at bounds (for documentation)
	MinCatchRate = 0.15 // 15% for strongest Pokémon
	MaxCatchRate = 0.95 // 95% for weakest Pokémon
)
