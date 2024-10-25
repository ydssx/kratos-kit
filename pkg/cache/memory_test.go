package cache

import (
	"testing"
	"time"
)

func TestMemoryCache_Set(t *testing.T) {
	c := NewMemoryCache()
	// Set with no expiration
	c.Set("key1", "value1", 0)

	v, ok := c.Get("key1")
	if !ok || v != "value1" {
		t.Error("Expected value to be set with no expiration")
	}

	// Set with expiration
	expiration := 1 * time.Second
	c.Set("key2", "value2", expiration)

	v, ok = c.Get("key2")
	if !ok || v != "value2" {
		t.Error("Expected value to be available before expiration")
	}

	// Wait for expiration
	time.Sleep(expiration + 100*time.Millisecond)

	_, ok = c.Get("key2")
	if ok {
		t.Error("Expected value to be unavailable after expiration")
	}

	// Update to reset expiration
	c.Set("key2", "newvalue2", expiration)

	v, ok = c.Get("key2")
	if !ok || v != "newvalue2" {
		t.Error("Expected updated value to reset expiration")
	}

	// Delete key
	c.Delete("key1")
	_, ok = c.Get("key1")
	if ok {
		t.Error("Expected key to be deleted")
	}
}
