package api

import (
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
)

// Create the struct to contain Pokemon Location Areas
type LocationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Area struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func ApiRequest(url string, cache *pokecache.Cache) ([]byte, error) {
	body, exists := cache.Get(url)
	if !exists {
		res, err := http.Get(url)
		if err != nil {
			return body, fmt.Errorf("Error requesting data: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return body, fmt.Errorf("Request failed with status code: %v", res.StatusCode)
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return body, fmt.Errorf("Error reading body: %w", err)
		}

		cache.Add(url, body)
	}

	return body, nil
}
