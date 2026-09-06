package store

import "sync"

// Activity describes a single Discord activity in the cache.
type Activity struct {
	Name    string
	Type    int
	Details string
	State   string
	SyncID  string
}

// Presence is a cached Discord presence.
type Presence struct {
	Status       string
	ClientStatus any
	Activities   []Activity
}

// Store is a concurrency-safe in-memory presence cache.
type Store struct {
	mu   sync.RWMutex
	data map[string]Presence
}

// New returns an empty Store.
func New() *Store {
	return &Store{data: make(map[string]Presence)}
}

// Set stores the presence for the given user ID.
func (s *Store) Set(id string, p Presence) {
	s.mu.Lock()
	s.data[id] = p
	s.mu.Unlock()
}

// Get returns the stored presence for the given user ID and whether it was found.
func (s *Store) Get(id string) (Presence, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.data[id]
	if ok {
		p.Activities = append([]Activity(nil), p.Activities...)
	}
	return p, ok
}
