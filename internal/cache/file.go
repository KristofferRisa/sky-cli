package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	// ErrCacheMiss is returned when a cache key is not found
	ErrCacheMiss = errors.New("cache miss")

	// ErrCacheExpired is returned when a cached value has expired
	ErrCacheExpired = errors.New("cache expired")
)

// cacheEntry represents a cached value with metadata
type cacheEntry struct {
	Value     []byte    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
}

// FileCache implements file-based caching
type FileCache struct {
	dir string
}

// NewFileCache creates a new file-based cache
func NewFileCache(dir string) (*FileCache, error) {
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &FileCache{
		dir: dir,
	}, nil
}

// Get retrieves a value from cache
func (c *FileCache) Get(key string) ([]byte, error) {
	filename := c.keyToFilename(key)

	// Read cache file
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	// Unmarshal entry
	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		// Invalid cache entry, delete it
		os.Remove(filename)
		return nil, ErrCacheMiss
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		// Clean up expired entry
		os.Remove(filename)
		return nil, ErrCacheExpired
	}

	return entry.Value, nil
}

// Set stores a value in cache with a TTL
func (c *FileCache) Set(key string, value []byte, ttl time.Duration) error {
	filename := c.keyToFilename(key)

	// Create entry
	entry := cacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}

	// Marshal entry
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// Delete removes a value from cache
func (c *FileCache) Delete(key string) error {
	filename := c.keyToFilename(key)

	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete cache file: %w", err)
	}

	return nil
}

// Clear removes all cached values
func (c *FileCache) Clear() error {
	// Read directory
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	// Delete all cache files
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := filepath.Join(c.dir, entry.Name())
		if err := os.Remove(filename); err != nil {
			return fmt.Errorf("failed to delete cache file %s: %w", entry.Name(), err)
		}
	}

	return nil
}

// Has checks if a key exists and is not expired
func (c *FileCache) Has(key string) bool {
	_, err := c.Get(key)
	return err == nil
}

// keyToFilename converts a cache key to a filename
func (c *FileCache) keyToFilename(key string) string {
	// Hash the key to create a safe filename
	hash := sha256.Sum256([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	return filepath.Join(c.dir, hashStr+".json")
}

// CleanExpired removes all expired cache entries
func (c *FileCache) CleanExpired() error {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := filepath.Join(c.dir, entry.Name())

		// Try to read and check expiration
		data, err := os.ReadFile(filename)
		if err != nil {
			continue
		}

		var cacheEntry cacheEntry
		if err := json.Unmarshal(data, &cacheEntry); err != nil {
			// Invalid entry, delete it
			os.Remove(filename)
			continue
		}

		// Delete if expired
		if time.Now().After(cacheEntry.ExpiresAt) {
			os.Remove(filename)
		}
	}

	return nil
}
