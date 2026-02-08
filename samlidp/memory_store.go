package samlidp

import (
	"encoding/json"
	"strings"
	"sync"
)

// DefaultMaxStoreItems is the default maximum number of items held by a
// MemoryStore before the oldest entries are evicted.
const DefaultMaxStoreItems = 10000

// MemoryStore is an implementation of Store that resides completely
// in memory. It has a configurable maximum capacity; when the limit
// is reached, the oldest entries are evicted.
type MemoryStore struct {
	mu       sync.RWMutex
	data     map[string]string
	keys     []string // insertion-ordered for eviction
	MaxItems int      // 0 means DefaultMaxStoreItems
}

func (s *MemoryStore) maxItems() int {
	if s.MaxItems > 0 {
		return s.MaxItems
	}
	return DefaultMaxStoreItems
}

// Get fetches the data stored in `key` and unmarshals it into `value`.
func (s *MemoryStore) Get(key string, value interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return ErrNotFound
	}
	return json.Unmarshal([]byte(v), value)
}

// Put marshals `value` and stores it in `key`.
func (s *MemoryStore) Put(key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data == nil {
		s.data = map[string]string{}
	}

	buf, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if _, exists := s.data[key]; !exists {
		s.keys = append(s.keys, key)
	}
	s.data[key] = string(buf)

	// Evict oldest entries if over capacity
	for len(s.data) > s.maxItems() && len(s.keys) > 0 {
		oldest := s.keys[0]
		s.keys = s.keys[1:]
		if _, ok := s.data[oldest]; ok {
			delete(s.data, oldest)
		}
	}
	return nil
}

// Delete removes `key`
func (s *MemoryStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	// Note: we don't remove from s.keys for efficiency; the eviction
	// loop in Put handles stale keys gracefully since they won't be
	// in s.data.
	return nil
}

// List returns all the keys that start with `prefix`. The prefix is
// stripped from each returned value. So if keys are ["aa", "ab", "cd"]
// then List("a") would produce []string{"a", "b"}
func (s *MemoryStore) List(prefix string) ([]string, error) {
	rv := []string{}
	for k := range s.data {
		if strings.HasPrefix(k, prefix) {
			rv = append(rv, strings.TrimPrefix(k, prefix))
		}
	}
	return rv, nil
}
