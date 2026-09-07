package store

import (
	"sync"
	"time"
)

// Discord activity type constants.
const (
	ActivityTypePlaying   = iota // Playing a game
	ActivityTypeStreaming        // Streaming
	ActivityTypeListening        // Listening to music
	ActivityTypeWatching         // Watching media
)

// Activity represents a single Discord activity in the cache.
type Activity struct {
	Name    string `json:"name"`
	Type    int    `json:"type"`
	Details string `json:"details"`
	State   string `json:"state"`
	URL     string `json:"url"`
	SyncID  string `json:"sync_id"`
}

// Presence is a cached Discord user presence.
type Presence struct {
	Status       string            `json:"status"`
	ClientStatus map[string]string `json:"client_status"`
	Activities   []Activity        `json:"activities"`
}

type entry struct {
	presence Presence
	updated  time.Time
}

// Store is a concurrency-safe in-memory presence cache.
type Store struct {
	mu   sync.RWMutex
	data map[string]entry
}

// New returns an empty Store. If ttl is positive, a background janitor
// removes presences that have not been updated within ttl.
func New(ttl time.Duration) *Store {
	s := &Store{data: make(map[string]entry)}
	if ttl > 0 {
		s.startJanitor(ttl)
	}
	return s
}

// Set stores the presence for the given user ID.
func (s *Store) Set(id string, p Presence) {
	s.mu.Lock()
	s.data[id] = entry{presence: p, updated: time.Now()}
	s.mu.Unlock()
}

// Get returns the stored presence for the given user ID and whether it was found.
func (s *Store) Get(id string) (Presence, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.data[id]
	if !ok {
		return Presence{}, false
	}
	e.presence.Activities = append([]Activity(nil), e.presence.Activities...)
	return e.presence, true
}

func (s *Store) startJanitor(ttl time.Duration) {
	interval := max(ttl/4, time.Minute)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			s.sweep(time.Now().Add(-ttl))
		}
	}()
}

func (s *Store) sweep(before time.Time) {
	s.mu.Lock()
	for id, e := range s.data {
		if e.updated.Before(before) {
			delete(s.data, id)
		}
	}
	s.mu.Unlock()
}
