// Package types provides the set of useful data structures.
package types

import (
	"sort"
	"sync"
)

// SafeMap provides a storage based on a map.
//
// It is safe to use in concurrent mode.
// The storage is protected by the mutex.
type SafeMap struct {
	mu      sync.Mutex // Protects storage below
	storage map[string]interface{}
}

// NewSafeMap returns a ready to use instance of SafeMap.
func NewSafeMap() *SafeMap {
	return &SafeMap{
		storage: make(map[string]interface{}),
	}
}

// Get returns object.
func (s *SafeMap) Get(key string) (interface{}, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.get(key)
}

// Set sets the object.
func (s *SafeMap) Set(key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.set(key, value)
}

// Del deletes the object.
func (s *SafeMap) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.del(key)
}

// Drain returns all elements as slice and removes keys from the storage.
// Elements are sorted by key.
func (s *SafeMap) Drain() []interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.drain()
}

// Len returns count of elements in the storage.
func (s *SafeMap) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.len()
}

// Keys returns all the keys as a slice.
func (s *SafeMap) Keys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.keys()
}

// get returns requested value and flag.
func (s *SafeMap) get(key string) (interface{}, bool) {
	r, ok := s.storage[key]

	return r, ok
}

// set sets value in the storage by key.
func (s *SafeMap) set(key string, value interface{}) error {
	s.storage[key] = value

	return nil
}

// del deletes the key from the storage.
func (s *SafeMap) del(key string) {
	delete(s.storage, key)
}

// drain returns values as a slice and removes data from the storage.
// values in the resulted slice are sorted by key.
func (s *SafeMap) drain() []interface{} {
	data := make([]interface{}, 0, len(s.storage))
	keys := make([]string, 0, len(s.storage))

	for k := range s.storage {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		data = append(data, s.storage[k])
		delete(s.storage, k)
	}

	return data
}

// len returns len of the storage.
func (s *SafeMap) len() int {
	return len(s.storage)
}

// keys returns a slice of keys.
func (s *SafeMap) keys() []string {
	keys := make([]string, 0, len(s.storage))

	for k := range s.storage {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
