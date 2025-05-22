package diskcache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Define a simple struct to test caching and encryption
type TestData struct {
	Name  string
	Value int
}

func TestCache(t *testing.T) {
	// Test 1: Test Set and Get with encryption and decryption
	t.Run("Set and Get data", func(t *testing.T) {
		// Create a simple struct to store in the cache
		testData := TestData{Name: "Test", Value: 123}

		// Store the test data in the cache with a 1-minute TTL
		err := Set("testKey", testData, 1*time.Minute)
		assert.NoError(t, err, "Expected no error while setting cache item")

		// Retrieve the data from the cache
		var result TestData
		found, err := Get("testKey", &result)
		assert.NoError(t, err, "Expected no error while getting cache item")
		assert.True(t, found, "Expected cache item to be found")
		assert.Equal(t, testData.Name, result.Name, "Expected name to match")
		assert.Equal(t, testData.Value, result.Value, "Expected value to match")
	})

	// Test 2: Check if expired item is not found
	t.Run("Cache item expiration", func(t *testing.T) {
		// Create a simple struct to store in the cache with a very short TTL
		testData := TestData{Name: "ExpiredTest", Value: 999}

		// Store the test data in the cache with a 1-second TTL
		err := Set("expiredKey", testData, 1*time.Second)
		assert.NoError(t, err, "Expected no error while setting cache item")

		// Wait for the cache item to expire
		time.Sleep(2 * time.Second)

		// Attempt to retrieve the expired data from the cache
		var result TestData
		found, err := Get("expiredKey", &result)
		assert.NoError(t, err, "Expected no error while getting expired cache item")
		assert.False(t, found, "Expected cache item to be expired and not found")
	})

	// Test 3: Check if key does not exist in cache
	t.Run("Key does not exist", func(t *testing.T) {
		// Try to retrieve a non-existent cache item
		var result TestData
		found, _ := Get("nonExistentKey", &result)
		assert.False(t, found, "Expected cache item to not be found")
	})
}
