package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pokedexcli/internal/pokecache"
	"strings"
	"testing"
	"time"
)

func TestApiRequestCacheMiss(t *testing.T) {
	// Create test data
	testData := locationArea{
		Count: 1281,
		Next:  "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
		Results: []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			{Name: "canalave-city-area", URL: "https://pokeapi.co/api/v2/location-area/1/"},
			{Name: "eterna-city-area", URL: "https://pokeapi.co/api/v2/location-area/2/"},
		},
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testData)
	}))
	defer server.Close()

	// Create cache and make request
	cache := pokecache.NewCache(5 * time.Minute)
	result, err := ApiRequest(server.URL, cache)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if result.Count != testData.Count {
		t.Errorf("expected count %d, got %d", testData.Count, result.Count)
		return
	}

	if len(result.Results) != len(testData.Results) {
		t.Errorf("expected %d results, got %d", len(testData.Results), len(result.Results))
		return
	}

	// Verify data was added to cache
	cachedData, exists := cache.Get(server.URL)
	if !exists {
		t.Errorf("expected data to be cached")
		return
	}

	if cachedData == nil {
		t.Errorf("cached data should not be nil")
		return
	}
}

func TestApiRequestCacheHit(t *testing.T) {
	// Create test data
	testData := locationArea{
		Count: 1281,
		Next:  "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
		Results: []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			{Name: "cached-area-1", URL: "https://pokeapi.co/api/v2/location-area/1/"},
			{Name: "cached-area-2", URL: "https://pokeapi.co/api/v2/location-area/2/"},
		},
	}

	// Create cache and pre-populate it
	cache := pokecache.NewCache(5 * time.Minute)
	testJSON, _ := json.Marshal(testData)
	testURL := "https://test-cached.example.com"
	cache.Add(testURL, testJSON)

	// Create server that should NOT be called (since we're testing cache hit)
	serverCallCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCallCount++
		t.Errorf("server was called, but should have been a cache hit")
	}))
	defer server.Close()

	// Make request - should use cached data, not hit server
	result, err := ApiRequest(testURL, cache)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if serverCallCount > 0 {
		t.Errorf("server was called %d times, expected 0 (cache hit)", serverCallCount)
		return
	}

	if result.Count != testData.Count {
		t.Errorf("expected count %d from cache, got %d", testData.Count, result.Count)
		return
	}

	if len(result.Results) != len(testData.Results) {
		t.Errorf("expected %d results from cache, got %d", len(testData.Results), len(result.Results))
		return
	}

	if result.Results[0].Name != "cached-area-1" {
		t.Errorf("expected cached data, got different data")
		return
	}
}

func TestApiRequestHTTPError(t *testing.T) {
	// Create server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cache := pokecache.NewCache(5 * time.Minute)
	_, err := ApiRequest(server.URL, cache)

	if err == nil {
		t.Errorf("expected error for 404 response, got nil")
		return
	}

	if !strings.Contains(err.Error(), "Request failed with status code: 404") {
		t.Errorf("expected status code error, got: %v", err)
		return
	}

	// Verify error response was not cached
	_, exists := cache.Get(server.URL)
	if exists {
		t.Errorf("error response should not be cached")
		return
	}
}

func TestApiRequestInvalidJSON(t *testing.T) {
	// Create server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json {"))
	}))
	defer server.Close()

	cache := pokecache.NewCache(5 * time.Minute)
	_, err := ApiRequest(server.URL, cache)

	if err == nil {
		t.Errorf("expected error for invalid JSON, got nil")
		return
	}

	if !strings.Contains(err.Error(), "Error unmarshalling JSON") {
		t.Errorf("expected JSON unmarshalling error, got: %v", err)
		return
	}

	// Verify invalid JSON was still cached (since the HTTP request succeeded)
	_, exists := cache.Get(server.URL)
	if !exists {
		t.Errorf("response should be cached even if JSON is invalid")
		return
	}
}

func TestApiRequestCachedInvalidJSON(t *testing.T) {
	// Test what happens when cached data is invalid JSON
	cache := pokecache.NewCache(5 * time.Minute)
	testURL := "https://test-invalid-cached.example.com"

	// Add invalid JSON to cache
	invalidJSON := []byte("invalid json data {")
	cache.Add(testURL, invalidJSON)

	// Try to use the cached invalid JSON
	_, err := ApiRequest(testURL, cache)

	if err == nil {
		t.Errorf("expected error for invalid cached JSON, got nil")
		return
	}

	if !strings.Contains(err.Error(), "Error unmarshalling JSON") {
		t.Errorf("expected JSON unmarshalling error from cache, got: %v", err)
		return
	}
}

func TestApiRequestNetworkError(t *testing.T) {
	// Test with invalid URL that will cause network error
	cache := pokecache.NewCache(5 * time.Minute)
	_, err := ApiRequest("http://nonexistent-domain-12345.invalid", cache)

	if err == nil {
		t.Errorf("expected network error, got nil")
		return
	}

	if !strings.Contains(err.Error(), "Error requesting data") {
		t.Errorf("expected network error, got: %v", err)
		return
	}
}
