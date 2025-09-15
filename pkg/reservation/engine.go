package reservation

import (
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
