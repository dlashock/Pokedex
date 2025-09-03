package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Create the struct to contain Pokemon Location Areas
type locationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func ApiRequest() (locationArea, error) {
	var data locationArea

	url := "https://pokeapi.co/api/v2/location-area/"
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

	//fmt.Printf("Response Body:\n%s", body)
	if err := json.Unmarshal(body, &data); err != nil {
		return data, fmt.Errorf("Error unmarshalling JSON: %w", err)
	}
	return data, nil
}
