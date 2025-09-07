package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
)

// Create the struct to contain Pokemon Location Areas
type locationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func ApiRequest(url string, cache *pokecache.Cache) (locationArea, error) {
	val, exists := cache.Get(url)
	var data locationArea
	if !exists {
		res, err := http.Get(url)
		if err != nil {
			return data, fmt.Errorf("Error requesting data: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return data, fmt.Errorf("Request failed with status code: %v", res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return data, fmt.Errorf("Error reading body: %w", err)
		}

		cache.Add(url, body)

		if err := json.Unmarshal(body, &data); err != nil {
			return data, fmt.Errorf("Error unmarshalling JSON: %w", err)
		}
	} else {
		if err := json.Unmarshal(val, &data); err != nil {
			return data, fmt.Errorf("Error unmarshalling JSON: %w", err)
		}
	}

	return data, nil
}
