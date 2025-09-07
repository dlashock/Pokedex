package pokecache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestCacheMiss(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)
	
	// Try to get a key that was never added
	val, ok := cache.Get("https://nonexistent.com")
	
	if ok {
		t.Errorf("expected cache miss, but found key")
		return
	}
	
	if val != nil {
		t.Errorf("expected nil value on cache miss, got %v", val)
		return
	}
}

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()
	const interval = 5 * time.Second
	const numGoroutines = 20
	
	cache := NewCache(interval)
	var wg sync.WaitGroup
	
	// Launch multiple goroutines that add and get data simultaneously
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			key := fmt.Sprintf("https://example.com/%d", id)
			value := []byte(fmt.Sprintf("testdata-%d", id))
			
			// Add data
			cache.Add(key, value)
			
			// Immediately try to get it back
			retrievedVal, ok := cache.Get(key)
			if !ok {
				t.Errorf("goroutine %d: expected to find key %s", id, key)
				return
			}
			
			if string(retrievedVal) != string(value) {
				t.Errorf("goroutine %d: expected %s, got %s", id, string(value), string(retrievedVal))
				return
			}
		}(i)
	}
	
	wg.Wait()
}

func TestMutexRaceCondition(t *testing.T) {
	t.Parallel()
	const interval = 5 * time.Second
	const numGoroutines = 50
	const numOperations = 100
	
	cache := NewCache(interval)
	var wg sync.WaitGroup
	
	// Multiple goroutines hammering the same key
	sharedKey := "https://shared.example.com"
	
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			for j := 0; j < numOperations; j++ {
				// Alternate between adding and getting
				if j%2 == 0 {
					value := []byte(fmt.Sprintf("data-%d-%d", id, j))
					cache.Add(sharedKey, value)
				} else {
					cache.Get(sharedKey)
				}
			}
		}(i)
	}
	
	wg.Wait()
	
	// If we get here without panicking or hanging, the mutex is working
	// Final verification that we can still use the cache
	testValue := []byte("final-test")
	cache.Add(sharedKey, testValue)
	
	retrievedVal, ok := cache.Get(sharedKey)
	if !ok {
		t.Errorf("expected to find key after race condition test")
		return
	}
	
	if string(retrievedVal) != string(testValue) {
		t.Errorf("cache corruption detected after race condition test")
		return
	}
}
