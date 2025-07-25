package storage

import "sync"

type Store struct {
	mu         sync.RWMutex
	shortToURL map[string]string
	URLToShort map[string]string
}

func NewStore() *Store {
	return &Store{
		shortToURL: make(map[string]string),
		URLToShort: make(map[string]string),
	}
}
func (s *Store) FindByURL(longURL string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	code, ok := s.URLToShort[longURL]
	return code, ok
}
func (s *Store) Save(code, longURL string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shortToURL[code] = longURL
	s.URLToShort[longURL] = code
}

func (s *Store) Get(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.shortToURL[code]
	return url, ok
}
