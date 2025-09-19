package reservation

import (
	"time"

	"github.com/mdryaan/resgate/internal/models"
)

func (e *Engine) sweepExpiredLocked() {
	now := time.Now()
	for id, r := range e.store.Reservations {
		if r.ExpiresAt != nil && now.After(*r.ExpiresAt) {
			if p, ok := e.store.Pools[r.PoolName]; ok {
				p.Release(r.Resources)
			}
			t := now
			r.Status = models.StatusExpired
			r.ReleasedAt = &t
			delete(e.store.Reservations, id)
			e.cache.Invalidate(id)
		}
	}
}

func (e *Engine) SweepExpired() []string {
	e.store.Lock()
	defer e.store.Unlock()
	var ids []string
	now := time.Now()
	for id, r := range e.store.Reservations {
		if r.ExpiresAt != nil && now.After(*r.ExpiresAt) {
			if p, ok := e.store.Pools[r.PoolName]; ok {
				p.Release(r.Resources)
			}
			t := now
			r.Status = models.StatusExpired
			r.ReleasedAt = &t
			delete(e.store.Reservations, id)
			e.cache.Invalidate(id)
			ids = append(ids, id)
		}
	}
	return ids
}
