package cache

import "github.com/mdryaan/resgate/internal/models"

type Cache interface {
	Get(id string) (*models.Reservation, bool)
	Set(id string, r *models.Reservation)
	Delete(id string)
	All() []*models.Reservation
	Invalidate(id string)
	Snapshot() map[string]*models.Reservation
	Len() int
}
