package reservation

import (
	"fmt"
	"time"

	"github.com/mdryaan/resgate/internal/models"
	"github.com/mdryaan/resgate/internal/utils"
)

func (e *Engine) Reserve(tenantName, poolName string, res models.Resources, ttl int) (*models.Reservation, error) {
	e.store.Lock()
	defer e.store.Unlock()

	if _, ok := e.store.Tenants[tenantName]; !ok {
		return nil, fmt.Errorf("tenant '%s' not found", tenantName)
	}

	p, ok := e.store.Pools[poolName]
	if !ok {
		return nil, fmt.Errorf("pool '%s' not found", poolName)
	}

	e.sweepExpiredLocked()

	for _, r := range e.store.Reservations {
		if r.TenantName == tenantName && r.PoolName == poolName && r.Status == models.StatusActive {
			return nil, fmt.Errorf("tenant '%s' already has an active reservation in pool '%s'", tenantName, poolName)
		}
	}

	if !p.CanFit(res) {
		avail := p.Available()
		return nil, fmt.Errorf("insufficient resources in '%s': need cpu=%d mem=%d gpu=%d, available cpu=%d mem=%d gpu=%d",
			poolName, res.CPU, res.Memory, res.GPU, avail.CPU, avail.Memory, avail.GPU)
	}

	p.Allocate(res)

	r := &models.Reservation{
		ID:         utils.NewID(),
		TenantName: tenantName,
		PoolName:   poolName,
		Resources:  res,
		Status:     models.StatusActive,
		TTL:        ttl,
		CreatedAt:  time.Now(),
	}

	if ttl > 0 {
		exp := time.Now().Add(time.Duration(ttl) * time.Second)
		r.ExpiresAt = &exp
	}

	e.store.Reservations[r.ID] = r
	e.store.History = append(e.store.History, r)
	e.cache.Set(r.ID, r)

	return r, nil
}

func (e *Engine) Unreserve(id string) error {
	e.store.Lock()
	defer e.store.Unlock()

	r, ok := e.store.Reservations[id]
	if !ok {
		return fmt.Errorf("reservation '%s' not found", id)
	}
	if r.Status != models.StatusActive {
		return fmt.Errorf("reservation '%s' is not active (status: %s)", id, r.Status)
	}

	if p, ok := e.store.Pools[r.PoolName]; ok {
		p.Release(r.Resources)
	}

	now := time.Now()
	r.Status = models.StatusReleased
	r.ReleasedAt = &now

	delete(e.store.Reservations, id)
	e.cache.Invalidate(id)

	return nil
}
