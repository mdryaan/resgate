package reservation

import (
	"sort"

	"github.com/mdryaan/resgate/internal/models"
	"github.com/mdryaan/resgate/internal/store"
	"github.com/mdryaan/resgate/pkg/cache"
	"github.com/mdryaan/resgate/pkg/pool"
	"github.com/mdryaan/resgate/pkg/tenant"
)

type Engine struct {
	store   *store.Store
	cache   *cache.Store
	pools   *pool.Manager
	tenants *tenant.Registry
}

func NewEngine(s *store.Store) *Engine {
	e := &Engine{
		store:   s,
		cache:   cache.NewStore(),
		pools:   pool.NewManager(s),
		tenants: tenant.NewRegistry(s),
	}
	for id, r := range s.Reservations {
		if r.IsActive() {
			e.cache.Set(id, r)
		}
	}
	return e
}

func (e *Engine) Save() error {
	return e.store.Save()
}

func (e *Engine) CreatePool(p *models.Pool) error              { return e.pools.Create(p) }
func (e *Engine) GetPool(name string) (*models.Pool, error)    { return e.pools.Get(name) }
func (e *Engine) ListPools() []*models.Pool                    { return e.pools.List() }

func (e *Engine) AddTenant(t *models.Tenant) error              { return e.tenants.Add(t) }
func (e *Engine) GetTenant(name string) (*models.Tenant, error) { return e.tenants.Get(name) }
func (e *Engine) ListTenants() []*models.Tenant                 { return e.tenants.List() }

func (e *Engine) ActiveReservations() []*models.Reservation {
	e.SweepExpired()
	e.store.RLock()
	defer e.store.RUnlock()
	result := make([]*models.Reservation, 0, len(e.store.Reservations))
	for _, r := range e.store.Reservations {
		result = append(result, r)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result
}

func (e *Engine) History() []*models.Reservation {
	e.store.RLock()
	defer e.store.RUnlock()
	out := make([]*models.Reservation, len(e.store.History))
	copy(out, e.store.History)
	return out
}
