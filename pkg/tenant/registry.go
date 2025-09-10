package tenant

import (
	"fmt"
	"sort"

	"github.com/mdryaan/resgate/internal/models"
	"github.com/mdryaan/resgate/internal/store"
)

type Registry struct {
	store *store.Store
}

func NewRegistry(s *store.Store) *Registry {
	return &Registry{store: s}
}

func (r *Registry) Add(t *models.Tenant) error {
	if t.Name == "" {
		return fmt.Errorf("tenant name cannot be empty")
	}
	if err := ValidatePriority(t.Priority); err != nil {
		return err
	}
	r.store.Lock()
	defer r.store.Unlock()
	if _, ok := r.store.Tenants[t.Name]; ok {
		return fmt.Errorf("tenant '%s' already exists", t.Name)
	}
	r.store.Tenants[t.Name] = t
	return nil
}

func (r *Registry) Get(name string) (*models.Tenant, error) {
	r.store.RLock()
	defer r.store.RUnlock()
	t, ok := r.store.Tenants[name]
	if !ok {
		return nil, fmt.Errorf("tenant '%s' not found", name)
	}
	return t, nil
}

func (r *Registry) List() []*models.Tenant {
	r.store.RLock()
	defer r.store.RUnlock()
	tenants := make([]*models.Tenant, 0, len(r.store.Tenants))
	for _, t := range r.store.Tenants {
		tenants = append(tenants, t)
	}
	sort.Slice(tenants, func(i, j int) bool {
		return tenants[i].Priority > tenants[j].Priority
	})
	return tenants
}
