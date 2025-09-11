package cache

import "github.com/mdryaan/resgate/internal/models"

func (s *Store) InvalidateAll() {
	s.items = make(map[string]*models.Reservation)
}

func InvalidateTenant(s *Store, tenantName string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var removed []string
	for id, r := range s.items {
		if r.TenantName == tenantName {
			delete(s.items, id)
			removed = append(removed, id)
		}
	}
	return removed
}

func InvalidatePool(s *Store, poolName string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var removed []string
	for id, r := range s.items {
		if r.PoolName == poolName {
			delete(s.items, id)
			removed = append(removed, id)
		}
	}
	return removed
}
