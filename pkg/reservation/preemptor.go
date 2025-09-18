package reservation

import (
	"fmt"
	"time"

	"github.com/mdryaan/resgate/internal/models"
)

func (e *Engine) Preempt(preemptorName, targetName string) ([]*models.Reservation, error) {
	e.store.Lock()
	defer e.store.Unlock()

	preemptor, ok := e.store.Tenants[preemptorName]
	if !ok {
		return nil, fmt.Errorf("tenant '%s' not found", preemptorName)
	}

	target, ok := e.store.Tenants[targetName]
	if !ok {
		return nil, fmt.Errorf("tenant '%s' not found", targetName)
	}

	if preemptor.Priority <= target.Priority {
		return nil, fmt.Errorf("cannot preempt: '%s' (priority %d) must outrank '%s' (priority %d)",
			preemptorName, preemptor.Priority, targetName, target.Priority)
	}

	e.sweepExpiredLocked()

	var preempted []*models.Reservation
	for id, r := range e.store.Reservations {
		if r.TenantName != targetName || r.Status != models.StatusActive {
			continue
		}
		if p, ok := e.store.Pools[r.PoolName]; ok {
			p.Release(r.Resources)
		}
		now := time.Now()
		r.Status = models.StatusPreempted
		r.ReleasedAt = &now
		preempted = append(preempted, r)
		delete(e.store.Reservations, id)
		e.cache.Invalidate(id)
	}

	if len(preempted) == 0 {
		return nil, fmt.Errorf("tenant '%s' has no active reservations to preempt", targetName)
	}

	return preempted, nil
}
