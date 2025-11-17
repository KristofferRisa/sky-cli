package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileCache(t *testing.T) {
	// Create temp directory for test
	tmpDir, err := os.MkdirTemp("", "sky-cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create cache
	cache, err := NewFileCache(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	t.Run("Set and Get", func(t *testing.T) {
		key := "test-key"
		value := []byte("test-value")
		ttl := 1 * time.Hour

		// Set value
		err := cache.Set(key, value, ttl)
		if err != nil {
			t.Fatalf("Set() failed: %v", err)
		}

		// Get value
		retrieved, err := cache.Get(key)
		if err != nil {
			t.Fatalf("Get() failed: %v", err)
		}

		if string(retrieved) != string(value) {
			t.Errorf("Get() = %s; want %s", string(retrieved), string(value))
		}
	})

	t.Run("Get non-existent key", func(t *testing.T) {
		_, err := cache.Get("non-existent-key")
		if err != ErrCacheMiss {
			t.Errorf("Get() error = %v; want %v", err, ErrCacheMiss)
		}
	})

	t.Run("Has method", func(t *testing.T) {
		key := "has-test-key"
		value := []byte("has-test-value")

		// Should not exist yet
		if cache.Has(key) {
			t.Error("Has() = true for non-existent key; want false")
		}

		// Set value
		cache.Set(key, value, 1*time.Hour)

		// Should exist now
		if !cache.Has(key) {
			t.Error("Has() = false for existing key; want true")
		}
	})

	t.Run("Expired entry", func(t *testing.T) {
		key := "expired-key"
		value := []byte("expired-value")
		ttl := 1 * time.Millisecond

		// Set with very short TTL
		err := cache.Set(key, value, ttl)
		if err != nil {
			t.Fatalf("Set() failed: %v", err)
		}

		// Wait for expiration
		time.Sleep(10 * time.Millisecond)

		// Should be expired
		_, err = cache.Get(key)
		if err != ErrCacheExpired {
			t.Errorf("Get() error = %v; want %v", err, ErrCacheExpired)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		key := "delete-test-key"
		value := []byte("delete-test-value")

		// Set value
		cache.Set(key, value, 1*time.Hour)

		// Verify it exists
		if !cache.Has(key) {
			t.Fatal("Has() = false; want true before delete")
		}

		// Delete
		err := cache.Delete(key)
		if err != nil {
			t.Fatalf("Delete() failed: %v", err)
		}

		// Verify it's gone
		if cache.Has(key) {
			t.Error("Has() = true after delete; want false")
		}
	})

	t.Run("Clear", func(t *testing.T) {
		// Set multiple values
		for i := 0; i < 5; i++ {
			key := filepath.Join("clear-test", string(rune('a'+i)))
			cache.Set(key, []byte("value"), 1*time.Hour)
		}

		// Clear all
		err := cache.Clear()
		if err != nil {
			t.Fatalf("Clear() failed: %v", err)
		}

		// Verify directory is empty (except .gitkeep if exists)
		entries, err := os.ReadDir(tmpDir)
		if err != nil {
			t.Fatalf("ReadDir() failed: %v", err)
		}

		for _, entry := range entries {
			if entry.Name() != ".gitkeep" {
				t.Errorf("Clear() left file: %s", entry.Name())
			}
		}
	})
}

func TestNoOpCache(t *testing.T) {
	cache := NewNoOpCache()

	t.Run("Get always returns cache miss", func(t *testing.T) {
		_, err := cache.Get("any-key")
		if err != ErrCacheMiss {
			t.Errorf("Get() error = %v; want %v", err, ErrCacheMiss)
		}
	})

	t.Run("Has always returns false", func(t *testing.T) {
		if cache.Has("any-key") {
			t.Error("Has() = true; want false")
		}
	})

	t.Run("Set does nothing", func(t *testing.T) {
		err := cache.Set("key", []byte("value"), 1*time.Hour)
		if err != nil {
			t.Errorf("Set() error = %v; want nil", err)
		}
	})
}
