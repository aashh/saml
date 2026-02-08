package samlidp

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestMemoryStoreEviction(t *testing.T) {
	store := MemoryStore{MaxItems: 3}

	// Fill to capacity
	assert.Check(t, store.Put("/a", "val-a"))
	assert.Check(t, store.Put("/b", "val-b"))
	assert.Check(t, store.Put("/c", "val-c"))

	var val string
	assert.Check(t, store.Get("/a", &val))
	assert.Check(t, is.Equal("val-a", val))

	// Adding a 4th item should evict the oldest (/a)
	assert.Check(t, store.Put("/d", "val-d"))

	err := store.Get("/a", &val)
	assert.Check(t, is.Equal(ErrNotFound, err), "/a should be evicted")

	assert.Check(t, store.Get("/b", &val))
	assert.Check(t, is.Equal("val-b", val))
	assert.Check(t, store.Get("/d", &val))
	assert.Check(t, is.Equal("val-d", val))
}

func TestMemoryStoreEvictionAfterDelete(t *testing.T) {
	store := MemoryStore{MaxItems: 3}

	assert.Check(t, store.Put("/a", "val-a"))
	assert.Check(t, store.Put("/b", "val-b"))
	assert.Check(t, store.Put("/c", "val-c"))

	// Delete /a, then add two more items
	assert.Check(t, store.Delete("/a"))
	assert.Check(t, store.Put("/d", "val-d"))
	// Store now has /b, /c, /d = 3 items (at capacity)

	assert.Check(t, store.Put("/e", "val-e"))
	// Eviction should skip the stale /a key and evict /b
	err := store.Get("/b", new(string))
	assert.Check(t, is.Equal(ErrNotFound, err), "/b should be evicted")

	assert.Check(t, store.Get("/c", new(string)))
	assert.Check(t, store.Get("/d", new(string)))
	assert.Check(t, store.Get("/e", new(string)))
}

func TestMemoryStoreUpdateDoesNotDuplicate(t *testing.T) {
	store := MemoryStore{MaxItems: 3}

	assert.Check(t, store.Put("/a", "val-a"))
	assert.Check(t, store.Put("/b", "val-b"))
	// Update /a (should not add a second /a to keys)
	assert.Check(t, store.Put("/a", "val-a-updated"))
	assert.Check(t, store.Put("/c", "val-c"))

	// All three should be present (no eviction since we updated, not added)
	var val string
	assert.Check(t, store.Get("/a", &val))
	assert.Check(t, is.Equal("val-a-updated", val))
	assert.Check(t, store.Get("/b", &val))
	assert.Check(t, store.Get("/c", &val))
}

func TestMemoryStoreDefaultMaxItems(t *testing.T) {
	store := MemoryStore{}
	assert.Check(t, is.Equal(DefaultMaxStoreItems, store.maxItems()))
}

func TestMemoryStoreManyEvictions(t *testing.T) {
	store := MemoryStore{MaxItems: 5}

	// Insert 20 items, only the last 5 should survive
	for i := 0; i < 20; i++ {
		assert.Check(t, store.Put(fmt.Sprintf("/item-%d", i), i))
	}

	// Items 0-14 should be evicted
	for i := 0; i < 15; i++ {
		err := store.Get(fmt.Sprintf("/item-%d", i), new(int))
		assert.Check(t, is.Equal(ErrNotFound, err),
			"item-%d should be evicted", i)
	}

	// Items 15-19 should be present
	for i := 15; i < 20; i++ {
		var val int
		assert.Check(t, store.Get(fmt.Sprintf("/item-%d", i), &val))
		assert.Check(t, is.Equal(i, val))
	}
}
