package pool

import (
	"fmt"
	"sort"

	"github.com/mdryaan/resgate/internal/models"
	"github.com/mdryaan/resgate/internal/store"
)

type Manager struct {
	store *store.Store
}

func NewManager(s *store.Store) *Manager {
	return &Manager{store: s}
}

func (m *Manager) Create(p *models.Pool) error {
	if err := ValidateNew(p); err != nil {
		return err
	}
	m.store.Lock()
	defer m.store.Unlock()
	if _, ok := m.store.Pools[p.Name]; ok {
		return fmt.Errorf("pool '%s' already exists", p.Name)
	}
	m.store.Pools[p.Name] = p
	return nil
}

func (m *Manager) Get(name string) (*models.Pool, error) {
	m.store.RLock()
	defer m.store.RUnlock()
	p, ok := m.store.Pools[name]
	if !ok {
		return nil, fmt.Errorf("pool '%s' not found", name)
	}
	return p, nil
}

func (m *Manager) List() []*models.Pool {
	m.store.RLock()
	defer m.store.RUnlock()
	pools := make([]*models.Pool, 0, len(m.store.Pools))
	for _, p := range m.store.Pools {
		pools = append(pools, p)
	}
	sort.Slice(pools, func(i, j int) bool {
		return pools[i].Name < pools[j].Name
	})
	return pools
}
