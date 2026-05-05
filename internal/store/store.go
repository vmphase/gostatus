package store

import "sync"

type Activity struct {
	Name    string
	Type    int
	Details string
	State   string
	SyncID  string
}

type Presence struct {
	Status       string
	ClientStatus any
	Activities   []Activity
}

type Store struct {
	mu   sync.RWMutex
	data map[string]Presence
}

func New() *Store {
	return &Store{data: make(map[string]Presence)}
}

func (s *Store) Set(id string, p Presence) {
	s.mu.Lock()
	s.data[id] = p
	s.mu.Unlock()
}

func (s *Store) Get(id string) (Presence, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.data[id]
	return p, ok
}
