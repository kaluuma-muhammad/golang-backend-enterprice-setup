package ratelimiter

import (
	"sync"
)

type MemoryStore struct {
	mu    sync.RWMutex
	store map[string]*Entry
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[string]*Entry),
	}
}

func (m *MemoryStore) Get(key string) (*Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.store[key]
	return entry, exists
}

func (m *MemoryStore) Set(key string, entry *Entry) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[key] = entry
}

func (m *MemoryStore) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.store, key)
}

func (m *MemoryStore) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]string, 0, len(m.store))
	for k := range m.store {
		keys = append(keys, k)
	}

	return keys
}
