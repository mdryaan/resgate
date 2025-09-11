package cache

import (
	"sort"
	"time"

	"github.com/mdryaan/resgate/internal/models"
)

type Snap struct {
	TakenAt      time.Time
	Reservations []*models.Reservation
}

func TakeSnapshot(s *Store) *Snap {
	snap := &Snap{
		TakenAt: time.Now(),
	}
	raw := s.All()
	sort.Slice(raw, func(i, j int) bool {
		return raw[i].CreatedAt.Before(raw[j].CreatedAt)
	})
	snap.Reservations = raw
	return snap
}
