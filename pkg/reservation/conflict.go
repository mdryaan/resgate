package reservation

import (
	"fmt"

	"github.com/mdryaan/resgate/internal/models"
)

func (e *Engine) detectConflictLocked(tenantName, poolName string) error {
	for _, r := range e.store.Reservations {
		if r.TenantName == tenantName && r.PoolName == poolName && r.Status == models.StatusActive {
			return fmt.Errorf("tenant '%s' already has active reservation '%s' in pool '%s'",
				tenantName, r.ID, poolName)
		}
	}
	return nil
}

func (e *Engine) HasConflict(tenantName, poolName string) bool {
	e.store.RLock()
	defer e.store.RUnlock()
	for _, r := range e.store.Reservations {
		if r.TenantName == tenantName && r.PoolName == poolName && r.IsActive() {
			return true
		}
	}
	return false
}
