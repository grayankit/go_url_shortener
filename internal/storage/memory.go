package storage

import "sync"

type Store struct {
	mu    sync.RWMutex
	links map[string]string
}

func NewStore() *Store {
	return &Store{
		links: make(map[string]string),
	}
}

func (s *Store) Save(code, longURL string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.links[code] = longURL
}

func (s *Store) Get(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.links[code]
	return url, ok
}
