package cache

import (
	"sync"

	"github.com/mdryaan/resgate/internal/models"
)

type Store struct {
	mu    sync.RWMutex
	items map[string]*models.Reservation
}

func NewStore() *Store {
	return &Store{items: make(map[string]*models.Reservation)}
}

func (s *Store) Get(id string) (*models.Reservation, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.items[id]
	return r, ok
}

func (s *Store) Set(id string, r *models.Reservation) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[id] = r
}

func (s *Store) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, id)
}

func (s *Store) All() []*models.Reservation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*models.Reservation, 0, len(s.items))
	for _, r := range s.items {
		out = append(out, r)
	}
	return out
}

func (s *Store) Invalidate(id string) { s.Delete(id) }

func (s *Store) Snapshot() map[string]*models.Reservation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap := make(map[string]*models.Reservation, len(s.items))
	for k, v := range s.items {
		snap[k] = v
	}
	return snap
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}
