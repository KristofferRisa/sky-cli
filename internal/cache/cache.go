package cache

import (
	"time"
)

// Cache is the interface for caching weather data
type Cache interface {
	// Get retrieves a value from cache
	Get(key string) ([]byte, error)

	// Set stores a value in cache with a TTL
	Set(key string, value []byte, ttl time.Duration) error

	// Delete removes a value from cache
	Delete(key string) error

	// Clear removes all cached values
	Clear() error

	// Has checks if a key exists and is not expired
	Has(key string) bool
}

// NoOpCache is a cache that does nothing (disabled cache)
type NoOpCache struct{}

// NewNoOpCache creates a new no-op cache
func NewNoOpCache() *NoOpCache {
	return &NoOpCache{}
}

// Get always returns an error (cache miss)
func (c *NoOpCache) Get(key string) ([]byte, error) {
	return nil, ErrCacheMiss
}

// Set does nothing
func (c *NoOpCache) Set(key string, value []byte, ttl time.Duration) error {
	return nil
}

// Delete does nothing
func (c *NoOpCache) Delete(key string) error {
	return nil
}

// Clear does nothing
func (c *NoOpCache) Clear() error {
	return nil
}

// Has always returns false
func (c *NoOpCache) Has(key string) bool {
	return false
}
