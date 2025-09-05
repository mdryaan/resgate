package store

import (
	"sync"

	"github.com/mdryaan/resgate/internal/models"
)

type Store struct {
	mu          sync.RWMutex
	Pools       map[string]*models.Pool
	Tenants     map[string]*models.Tenant
	Reservations map[string]*models.Reservation
	History     []*models.Reservation
	persistPath string
}

func New(path string) *Store {
	return &Store{
		Pools:        make(map[string]*models.Pool),
		Tenants:      make(map[string]*models.Tenant),
		Reservations: make(map[string]*models.Reservation),
		History:      []*models.Reservation{},
		persistPath:  path,
	}
}

func (s *Store) RLock()   { s.mu.RLock() }
func (s *Store) RUnlock() { s.mu.RUnlock() }
func (s *Store) Lock()    { s.mu.Lock() }
func (s *Store) Unlock()  { s.mu.Unlock() }
