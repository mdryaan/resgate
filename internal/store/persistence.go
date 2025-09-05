package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/mdryaan/resgate/internal/models"
)

type persistState struct {
	Pools        map[string]*models.Pool        `json:"pools"`
	Tenants      map[string]*models.Tenant      `json:"tenants"`
	Reservations map[string]*models.Reservation `json:"reservations"`
	History      []*models.Reservation          `json:"history"`
}

func (s *Store) Load() error {
	if s.persistPath == "" {
		return nil
	}

	data, err := os.ReadFile(s.persistPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	var state persistState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	if state.Pools != nil {
		s.Pools = state.Pools
	}
	if state.Tenants != nil {
		s.Tenants = state.Tenants
	}
	if state.Reservations != nil {
		s.Reservations = state.Reservations
	}
	if state.History != nil {
		s.History = state.History
	}

	return nil
}

func (s *Store) Save() error {
	if s.persistPath == "" {
		return nil
	}

	s.mu.RLock()
	state := persistState{
		Pools:        s.Pools,
		Tenants:      s.Tenants,
		Reservations: s.Reservations,
		History:      s.History,
	}
	s.mu.RUnlock()

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(s.persistPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(s.persistPath, data, 0644)
}
